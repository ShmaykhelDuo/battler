const currencies = {
    1: "white",
    2: "blue",
    3: "yellow",
    4: "purple",
    5: "star"
}

function balance() {
    return {
        balance: {},

        init() {
            this.fetch();
        },

        async fetch() {
            const response = await fetch("/money/balance");
            if (!response.ok) {
                if (response.status == 401) {
                    window.location.href = "/web/auth/signin"
                }
                return;
            }

            const result = await response.json();
            this.update(result);
        },

        update(json) {
            for (let id in currencies) {
                if (json[id]) {
                    this.balance[id] = json[id]
                } else {
                    this.balance[id] = 0
                }
            }
        }
    }
}
