let socket = null

function buildSkillButtonSvg(parent, skillNum, text) {
    const svg = document.createElementNS("http://www.w3.org/2000/svg", "svg")
    parent.replaceChildren(svg)

    const words = text.split(" ")
    let angle = 0
    if (skillNum == 0) {
        angle = 37.5
    } else if (skillNum == 3) {
        angle = -37.5
    }

    const step = 50 / Math.cos(angle * Math.PI / 180)

    for (let i = 0; i < words.length; i++) {
        const y = step * (i - (words.length - 1) / 2)

        const element = document.createElementNS("http://www.w3.org/2000/svg", "text")
        element.setAttributeNS(null, "y", y)
        element.setAttributeNS(null, "font-size", "50px")
        element.setAttributeNS(null, "text-anchor", "middle")
        element.setAttributeNS(null, "dominant-baseline", "middle")
        element.setAttributeNS(null, "transform", `rotate(${angle},0,${y})`)
        element.textContent = words[i]

        svg.appendChild(element)
    }

    const bbox = svg.getBBox()
    svg.setAttribute("viewBox", `${bbox.x} ${bbox.y} ${bbox.width} ${bbox.height}`)
}

function setCharacter(panel, state, character) {
    panel.getElementsByClassName("character-name")[0].innerText = `${character.name} (${character.number})`
    panel.getElementsByClassName("character-hp")[0].innerText = `${state.hp}/${state.max_hp} (${(state.hp / state.max_hp * 100).toFixed(2)}%)`
    const effectsElement = panel.getElementsByClassName("character-effects")[0]
    effectsElement.innerHTML = ""
    for (const effect of state.effects) {
        effectsElement.innerHTML += `
            <div class="effect">
                <div class="effect-name">${effect.type}</div>
            </div>
        `
    }

    for (let i = 0; i < 4; i++) {
        const element = panel.getElementsByClassName(`skill-${i}`)[0]
        element.style.backgroundColor = character.skills[i].colour
        buildSkillButtonSvg(element, i, character.skills[i].name)
    }
}

function setSkillAvailability(state) {
    for (let i = 0; i < 4; i++) {
        const disabled = !state.player_turn || state.as_opp || !state.character.skill_availabilities[i]
        const element = document.getElementById(`character-skill-${i}`)
        if (disabled) {
            element.setAttribute("disabled", "disabled")
        } else {
            element.removeAttribute("disabled")
        }
    }

    for (let i = 0; i < 4; i++) {
        const disabled = !state.player_turn || !state.as_opp || !state.opponent.skill_availabilities[i]
        const element = document.getElementById(`opponent-skill-${i}`)
        if (disabled) {
            element.setAttribute("disabled", "disabled")
        } else {
            element.removeAttribute("disabled")
        }
    }
}

function setSkillLog(state) {
    const isCharactersTurn = state.player_turn != state.as_opp
    const isGoingFirst = isCharactersTurn == state.turn.first
    let firstClass = "skill-log-character"
    let secondClass = "skill-log-opponent"
    let firstCharacter = characters[state.character.number]
    let secondCharacter = characters[state.opponent.number]
    if (!isGoingFirst) {
        [firstClass, secondClass] = [secondClass, firstClass]
        [firstCharacter, secondCharacter] = [secondCharacter, firstCharacter]
    }

    const skillLogElement = document.getElementById("skill-log")
    skillLogElement.innerHTML = ""
    for (const entry of state.skill_log.turns) {
        if (entry.first) {
            for (const skillNum of entry.first) {
                const char = "QWER"[skillNum]
                const skill = firstCharacter.skills[skillNum]
    
                skillLogElement.innerHTML += `
                    <div class="${firstClass}" style="color: ${skill.colour};">[${char}] ${skill.name}</div>
                `
            }
        }

        if (entry.second) {
            for (const skillNum of entry.second) {
                const char = "QWER"[skillNum]
                const skill = secondCharacter.skills[skillNum]
    
                skillLogElement.innerHTML += `
                    <div class="${secondClass}" style="color: ${skill.colour};">[${char}] ${skill.name}</div>
                `
            }
        }
    }
}

function showState(state) {
    const element = document.getElementById("state")
    element.innerText = JSON.stringify(state)

    const character = characters[state.character.number]
    const characterPanel = document.getElementById("character-panel")
    setCharacter(characterPanel, state.character, character)

    const opponent = characters[state.opponent.number]
    const opponentPanel = document.getElementById("opponent-panel")
    setCharacter(opponentPanel, state.opponent, opponent)

    setSkillAvailability(state)

    document.getElementById("turn-number").innerText = `Turn ${2 * state.turn.number - state.turn.first} [${state.turn.number}]`
    setSkillLog(state)
}

function battle() {
    const params = new URLSearchParams(window.location.search);

    const main = parseInt(params.get("main"));
    const secondary = parseInt(params.get("secondary"));

    socket = new WebSocket("/game/match");
    socket.addEventListener("open", (event) => {
        socket.send(JSON.stringify({
            type: 1,
            payload: {
                main: main,
                secondary: secondary
            }
        }));
    });

    socket.addEventListener("message", (event) => {
        const message = JSON.parse(event.data);

        switch (message.type) {
            case 4:
                showState(message.payload.state);
                break;
        
            default:
                break;
        }
    })
}

function skill(i) {
    socket.send(JSON.stringify({
        type: 5,
        payload: {
            move: i
        }
    }));
}
