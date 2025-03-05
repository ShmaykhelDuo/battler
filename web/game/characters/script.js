let playerCharacters = [];
let displayed = 0;

let main = 0;
let secondary = 0;

function characterNode(character) {
    const element = document.createElement("button");
    element.classList.add("character-button")
    element.onclick = () => showCharacter(character.number)
    element.style.backgroundColor = character.colour
    element.innerText = character.name

    return element;
}

function displayCharacters() {
    const nodes = playerCharacters.map(characterNode);

    const container = document.getElementById("characters");
    container.replaceChildren(...nodes);

    showCharacter(playerCharacters[0].number)
}

async function fetchCharacters() {
    const response = await fetch("/game/characters");
    if (!response.ok) {
        if (response.status == 401) {
            window.location.href = "/web/auth/signin"
        }
        handleError(response);
        return;
    }

    const result = await response.json();
    playerCharacters = result.map(c => characters[c.number]);
    displayCharacters();
}

function showCharacter(number) {
    const character = characters[number];

    document.getElementById("character-name").innerText = `${character.name} (${character.number})`
    document.getElementById("character-rarity").innerText = character.rarity
    document.getElementById("character-level").innerText = undefined
    document.getElementById("character-matches").innerText = undefined
    document.getElementById("character-matches-won").innerText = undefined
    document.getElementById("character-tags").innerText = character.tags.join(", ")
    document.getElementById("character-colours").innerText = character.skills.map(skill => skill.colour).join(", ")

    const skills = document.getElementById("character-skills")
    skills.innerHTML = ""
    for (const skill of character.skills) {
        skills.innerHTML += `
            <div class="skill" style="background-color: ${skill.colour}">${skill.name}</div>
        `
    }

    displayed = number;
}

function setAsMain() {
    if (displayed == secondary) {
        secondary = 0
    }
    main = displayed;
    updateSelections()
}

function setAsSecondary() {
    if (displayed == main) {
        main = 0
    }
    secondary = displayed;
    updateSelections()
}

function clearSelection() {
    main = 0;
    secondary = 0;
    updateSelections()
}

function chooseRandom() {
    updateSelections()
}

function updateSelections() {
    const element = document.getElementById("battle-button")
    if (main != 0 && secondary != 0) {
        element.removeAttribute("disabled")
    } else {
        element.setAttribute("disabled", "disabled")
    }

    const mainElement = document.getElementById("selected-main")
    updateSelectedItem(mainElement, main)

    const secondaryElement = document.getElementById("selected-secondary")
    updateSelectedItem(secondaryElement, secondary)
}

function updateSelectedItem(element, selection) {
    if (selection == 0) {
        element.classList.add("selected-item-empty")
        element.innerText = ""
        element.style = ""
    } else {
        const character = characters[selection]
        element.classList.remove("selected-item-empty")
        element.innerText = character.name
        element.style.backgroundColor = character.colour
    }
}

function battle() {
    window.location.href = `/web/game/battle?main=${main}&secondary=${secondary}`;
}
