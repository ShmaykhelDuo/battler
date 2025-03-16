function characterSelection() {
    return {
        characters: characters,
        availableCharacters: [],
        displayed: null,
        main: 0,
        secondary: 0,

        init() {
            this.fetchCharacters()
        },

        async fetchCharacters() {
            const response = await fetch("/game/characters");
            if (!response.ok) {
                if (response.status == 401) {
                    window.location.href = "/web/auth/signin"
                }
                return;
            }
        
            const result = await response.json();
            this.availableCharacters = result.map(c => characters[c.number]);
            this.displayed = this.availableCharacters[0];
        },

        displayCharacter(number) {
            this.displayed = characters[number];
        },

        setAsMain() {
            if (this.displayed.number == this.secondary) {
                this.secondary = 0
            }
            this.main = this.displayed.number;
        },

        setAsSecondary() {
            if (this.displayed.number == this.main) {
                this.main = 0
            }
            this.secondary = this.displayed.number;
        },

        clearSelection() {
            this.main = 0;
            this.secondary = 0;
        },

        chooseRandom() {
            const n = this.availableCharacters.length;
            let first = Math.floor(Math.random() * n);
            let second = Math.floor(Math.random() * (n - 1));
            if (second >= first) {
                second++;
            }
            this.main = this.availableCharacters[first].number;
            this.secondary = this.availableCharacters[second].number;
        },

        battle() {
            window.location.href = `/web/game/battle?main=${this.main}&secondary=${this.secondary}`;
        }
    }
}
