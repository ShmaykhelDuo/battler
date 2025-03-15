function conversion() {
    return {
        balance: balance(),
        conversion: null,
        selected_currency: 1,
        amount: 1,
        seconds_left: 0,
        progress: 0.0,

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

        timeLeft() {
            const now = new Date().getTime();
            return new Date(this.conversion.finishes_at).getTime() - now;
        },

        totalDuration() {
            return new Date(this.conversion.finishes_at).getTime() - new Date(this.conversion.started_at).getTime();
        },

        countdown() {
            const timeLeft = this.timeLeft();
            if (timeLeft < 0) {
                this.seconds_left = 0;
                this.progress = 100.0;
                return
            }

            this.seconds_left = Math.ceil(timeLeft / 1000);
            const totalDuration = this.totalDuration();
            this.progress = 100 * (totalDuration - timeLeft) / totalDuration;

            let x = setInterval(() => {
                const timeLeft = this.timeLeft();

                if (timeLeft < 0) {
                    clearInterval(x);
                    return;
                }

                this.seconds_left = Math.ceil(timeLeft / 1000);
                const totalDuration = this.totalDuration();
                this.progress = 100 * (totalDuration - timeLeft) / totalDuration;
            }, 10);
        },
    }
}
