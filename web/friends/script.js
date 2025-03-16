function friends() {
    return {
        friends: [],
        outgoingRequests: [],
        incomingRequests: [],

        init() {
            this.fetchAll()
        },

        fetchAll() {
            this.fetchFriends()
            this.fetchOutgoingRequests()
            this.fetchIncomingRequests()
        },

        async fetchFriends() {
            this.friends = this.fetchUrl("/friends")
        },

        async fetchOutgoingRequests() {
            this.outgoingRequests = this.fetchUrl("/friends/outgoing")
        },

        async fetchIncomingRequests() {
            this.incomingRequests = this.fetchUrl("/friends/incoming")
        },

        async fetchUrl(url) {
            const response = await fetch(url);

            if (!response.ok) {
                if (response.status == 401) {
                    window.location.href = "/web/auth/signin";
                }
                return null
            }

            return await response.json()
        },

        async addFriendLink(id) {
            const response = await fetch(`/friends/${id}`, {
                method: "POST"
            })

            if (!response.ok) {
                if (response.status == 401) {
                    window.location.href = "/web/auth/signin";
                }
                return
            }

            this.fetchAll()
        },

        async removeFriendLink(id) {
            const response = await fetch(`/friends/${id}`, {
                method: "DELETE"
            })

            if (!response.ok) {
                if (response.status == 401) {
                    window.location.href = "/web/auth/signin";
                }
                return
            }

            this.fetchAll()
        }
    }
}