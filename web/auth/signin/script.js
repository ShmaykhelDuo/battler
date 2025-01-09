async function handleError(response) {
    const json = await response.json();
    const message = parseError(json.error);

    document.getElementById("error").innerText = message;
}

async function signin() {
    const username = document.getElementById("username").value;
    const password = document.getElementById("password").value;

    const response = await fetch("/auth/signin", {
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
