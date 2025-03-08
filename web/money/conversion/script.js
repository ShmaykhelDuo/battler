function conversion() {
    return {
        balance: balance(),
        conversion: null,
        selected_currency: 1,
        amount: 1,
        time_left: 0,

        init() {
            this.balance.init();
            this.fetch();
        },

        async fetch() {
            const response = await fetch("/money/conversion");

            if (!response.ok) {
                if (response.status == 401) {
                    window.location.href = "/web/auth/signin";
                }
                else if (response.status == 404) {
                    this.conversion = null;
                }
                return
            }

            this.conversion = await response.json()
            this.countdown()
        },

        async convert() {
            const response = await fetch("/money/conversion", {
                method: "POST",
                body: JSON.stringify({
                    currency_id: this.selected_currency,
                    amount: this.amount
                })
            })

            if (!response.ok) {
                if (response.status == 401) {
                    window.location.href = "/web/auth/signin";
                }
                return
            }

            this.balance.fetch()

            this.conversion = await response.json()
            this.countdown()
        },

        async claim() {
            const response = await fetch("/money/conversion/claim", {
                method: "POST"
            });
            if (!response.ok) {
                return;
            }
        
            this.conversion = null;
            this.balance.update(await response.json())
        },

        secondsLeft() {
            const now = new Date().getTime();
            const distance = new Date(this.conversion.finishes_at).getTime() - now;
            return Math.ceil(distance / 1000);
        },

        countdown() {
            this.time_left = this.secondsLeft();
            if (this.time_left < 0) {
                this.time_left = 0;
                return
            }

            let x = setInterval(() => {
                this.time_left = this.secondsLeft();
        
                if (this.time_left <= 0) {
                    clearInterval(x);
                }
            }, 1000);
        },
    }
}
