function signinData() {
    return {
        username: null,
        password: null,
        error: null,

        async signin() {
            const response = await fetch("/auth/signin", {
                method: "POST",
                body: JSON.stringify({
                    username: this.username,
                    password: this.password
                })
            });
            if (!response.ok) {
                this.error = handleError(response)
                return;
            }
        
            window.location.href = "/web/";
        }
    }
}
