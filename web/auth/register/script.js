async function handleError(response) {
    const json = await response.json();
    const message = parseError(json.error);

    document.getElementById("error").innerText = message;
}

async function register() {
    const username = document.getElementById("username").value;
    const password = document.getElementById("password").value;
    const repeatedPassword = document.getElementById("repeated_password").value;

    if (password != repeatedPassword) {
        document.getElementById("error").innerText = "Passwords do not match";
        return;
    }

    const response = await fetch("/auth/register", {
        method: "POST",
        body: JSON.stringify({
            username: username,
            password: password
        })
    });
    if (!response.ok) {
        handleError(response);
        return;
    }

    window.location.href = "/web/";
}
