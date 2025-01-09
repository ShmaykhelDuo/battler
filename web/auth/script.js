async function signout() {
    const response = await fetch("/auth/signout", { method: "POST" });
    if (!response.ok) {
        handleError(response);
        return;
    }

    window.location.href = "/web/auth/signin";
}