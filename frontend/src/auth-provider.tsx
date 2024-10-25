async function fetchUserInfo() {
  const response = await fetch("/oauth2/userinfo");
  const userInfo = await response.json();
  return userInfo;
}

const authProvider = {
  checkAuth: () => Promise.resolve(),
  checkError: () => Promise.resolve(),
  getIdentity: () => fetchUserInfo(),
  getPermissions: () => Promise.resolve(),
  logout: () => Promise.resolve(),
};

export default authProvider;
