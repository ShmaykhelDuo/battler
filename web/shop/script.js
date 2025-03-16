function chests() {
    return {
        balance: balance(),
        chests: [],

        init() {
            this.balance.init()
            this.fetch()
        },

        fetch() {
            fetch("/shop/chests")
                .then((response) => response.json())
                .then((json) => this.chests = json)
        },

        async buy(id) {
            const response = await fetch(`/shop/chests/${id}`, {
                method: "POST"
            });
            if (!response.ok) {
                if (response.status == 401) {
                    window.location.href = "/web/auth/signin"
                }
                return;
            }

            this.balance.fetch();
        }
    }
}
