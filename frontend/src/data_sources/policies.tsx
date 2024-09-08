import {
  createPolicy,
  deletePolicy,
  fetchPolicies,
  fetchPolicy,
} from "../api/policies";

import { castArray } from "lodash";

function normalizePolicy(policy) {
  policy.policy["odrl:permission"] = castArray(
    policy.policy["odrl:permission"]
  );
  for (const permission of policy.policy["odrl:permission"]) {
    if (permission["odrl:constraint"] !== undefined) {
      permission["odrl:constraint"] = castArray(permission["odrl:constraint"]);
    }
  }

  policy.policy["odrl:prohibition"] = castArray(
    policy.policy["odrl:prohibition"]
  );
  for (const prohibition of policy.policy["odrl:prohibition"]) {
    if (prohibition["odrl:constraint"] !== undefined)
      prohibition["odrl:constraint"] = castArray(
        prohibition["odrl:constraint"]
      );
  }

  policy.policy["odrl:obligation"] = castArray(
    policy.policy["odrl:obligation"]
  );
  for (const obligation of policy.policy["odrl:obligation"]) {
    if (obligation["odrl:constraint"] !== undefined)
      obligation["odrl:constraint"] = castArray(obligation["odrl:constraint"]);
  }

  return policy;
}

export async function getList(params) {
  const policies = await fetchPolicies(params.pagination);

  for (let i = 0; i < policies.length; i++) {
    policies[i] = normalizePolicy(policies[i]);
  }

  return {
    data: policies.map((policy: any) => ({
      ...policy,
      id: policy["@id"],
    })),
    pageInfo: {
      hasNextPage: policies.length === params.pagination.perPage,
      hasPreviousPage: params.pagination.page > 1,
    },
  };
}

export async function getOne(params) {
  const policy = await fetchPolicy(params.id);
  return {
    data: {
      ...normalizePolicy(policy),
      id: policy["@id"],
    },
  };
}

export async function remove(params) {
  await deletePolicy(params.id);
  return {
    data: {
      id: params.id,
    },
  };
}

export async function create(params) {
  const { permissions, ...data } = params.data;
  const policy = await createPolicy({
    ...data,
    policy: {
      ...params.data.policy,
      "@type": "odrl:Set",
      "odrl:permission": permissions.map((permission) => ({
        "odrl:action": {
          "@id": "odrl:use",
        },
        "odrl:constraint": [
          {
            "@type": "AtomicConstraint",
            "odrl:leftOperand": {
              "@id": permission.leftOperand,
            },
            "odrl:operator": {
              "@id": permission.operator,
            },
            "odrl:rightOperand": permission.rightOperand,
          },
        ],
      })),
    },
    "@context": {
      "@vocab": "https://w3id.org/edc/v0.0.1/ns/",
      edc: "https://w3id.org/edc/v0.0.1/ns/",
      odrl: "http://www.w3.org/ns/odrl/2/",
    },
  });
  return {
    data: {
      ...policy,
      id: policy["@id"],
    },
  };
}

export async function getMany(params) {
  const policies = await Promise.all(
    params.ids.map((id: any) => fetchPolicy(id))
  );

  for (let i = 0; i < policies.length; i++) {
    policies[i] = normalizePolicy(policies[i]);
  }

  return {
    data: policies.map((policy: any) => ({
      ...policy,
      id: policy["@id"],
    })),
  };
}
