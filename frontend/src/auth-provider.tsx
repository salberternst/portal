const authProvider = {
    checkAuth: () => Promise.resolve(),
    checkError:  () => {
        return Promise.resolve();
    },
    getIdentity: async () => {
        const response = await fetch(
            '/user/info',
        );
        const userInfo = await response.json()
        return { 
            id: userInfo.name,
            fullName: userInfo.name
        }
    },
    getPermissions: () => Promise.resolve(''),
};

export default authProvider;