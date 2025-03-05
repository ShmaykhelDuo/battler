const kinds = [
    "Internal",
    "Invalid request",
    "Invalid argument",
    "Not found",
    "Already exists",
    "Invalid username or password",
    "Not authenticated"
];

function parseError(error) {
    if (error.message) {
        return kinds[error.id] + ": " + error.message;
    }

    return kinds[error.id];
}

async function handleError(response) {
    const json = await response.json();
    const message = parseError(json.error);

    document.getElementById("error").innerText = message;
}
