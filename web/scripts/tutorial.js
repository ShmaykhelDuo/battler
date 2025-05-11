const tutorialScenario = [
    {
        page: "/web/",
        text: "Welcome to Battler! This is a tutorial designed to showcase the game and all its features. You can skip it by pressing 'Close' at any point. Press 'Next' to continue.",
    },
    {
        page: "/web/",
        text: "This is a welcome screen. It contains some info about your profile.",
    },
    {
        page: "/web/",
        text: "This screen, as well as most of the rest of them, contains info about your balance. It is placed on the top left corner of the screen.",
    },
    {
        page: "/web/",
        text: "There is a navigation panel on the top that allows navigating to other pages. Try to go to the character selection page using 'Choose chars' link.",
        nextPage: "/web/game/characters/"
    },
    {
        page: "/web/game/characters/",
        text: "This is a character selection screen. It displays the characters which are unlocked",
    },
    {
        page: "/web/game/characters/",
        text: "On the left side, the list of characters is situated. Try pressing on a character button to display info about them.",
        trigger: "characterSelectionDisplay"
    },
    {
        page: "/web/game/characters/",
        text: "On the right to the list, features of the character are displayed, such as name, number, description, and skills.",
    },
    {
        page: "/web/game/characters/",
        text: "You can find out more information about each skill by hovering over its panel.",
    },
    {
        page: "/web/game/characters/",
        text: "But the main purpose of this screen is not just to learn about the characters. It is to select them for a battle.",
    },
    {
        page: "/web/game/characters/",
        text: "To participate in a match you need to select two characters: main and secondary. This is required so the matchmaking service is able to resolve conflicts between players who selected the same characters.",
    },
    {
        page: "/web/game/characters/",
        text: "The panel on the bottom serves this purpose. You can select characters manually using 'Set as Main' and 'Set as Secondary' buttons, as well as select random characters using the 'Choose random' button. Select the characters using any method.",
        trigger: "characterSelectionBothSelected"
    },
    {
        page: "/web/game/characters/",
        text: "Now, when the characters are selected, you can request a match using the 'Battle' button. Let's go!",
        nextPage: "/web/game/match/"
    },
    {
        page: "/web/game/match/",
        text: "Welcome to the match screen! This is the screen which hopefully you will spend the most time in.",
    },
    {
        page: "/web/game/match/",
        text: "Let's have a quick walkthrough of the match screen interface. There are two similar-looking panels on the left and on the right. These represent your (left) and your opponent's (right) characters.",
    },
    {
        page: "/web/game/match/",
        text: "Each of these panels is divided into multiple sections. Character information, like number and HP, is placed on the top.",
    },
    {
        page: "/web/game/match/",
        text: "Your main objective in this game is to decrease your opponent's HP while maintaining your own. If your opponent's HP reach zero, you immediately win and vice versa. However, if by the end of the game neither character's HP reach zero, winner is the player which character's HP value is larger.",
    },
    {
        page: "/web/game/match/",
        text: "In order to achieve your goal, you need to use the skills provided by your character. There are four buttons on the bottom for this purpose, each representing a skill. Like on the character selection screen, you can read skill's description by hovering over its button.",
    },
    {
        page: "/web/game/match/",
        text: "When using skills that damage the opponent, note that any damage dealt of colour is decreased by the number specified in the defences of that colour. This information can be viewed by hovering over character's HP.",
    },
    {
        page: "/web/game/match/",
        text: "Some skills can apply an effect on a character. These effects are displayed on the left side of character's panel, just below their HP. More information on each effect is displayed on hover.",
    },
    {
        page: "/web/game/match/",
        text: "Now, when you know the basics of the game mechanics, let's talk more about the match pacing. During the match, characters take turns, but its length is limited to 10 turns per character.",
    },
    {
        page: "/web/game/match/",
        text: "Information about the current turn, as well as the history of used skills, is displayed on the middle panel.",
    },
    {
        page: "/web/game/match/",
        text: "Good luck! The tutorial will continue when the match is finished.",
        nextPage: "/web/game/rewards/"
    },
    {
        page: "/web/game/rewards/",
        text: "After each match you receive white dust and experience with the possibility of level-up.",
    },
    {
        page: "/web/game/rewards/",
        text: "The dust you collect can be converted into another types. Go to the conversion page by clicking on 'Conversion' link on the navigation panel.",
        nextPage: "/web/money/conversion/"
    },
    {
        page: "/web/money/conversion/",
        text: "Each type of dust can be converted into the next one at a specified rate. This is a one way process. That means since you converted white dust to blue dust, you can't get the white dust back from it.",
    },
    {
        page: "/web/money/conversion/",
        text: "Each conversion takes some amount of time depending on the dust type and the amount. At any point of time, you can do only one conversion.",
    },
    {
        page: "/web/money/conversion/",
        text: "Having discussed the theory, let's try to convert some dust, to be precise, 2 white dust to 1 blue dust.",
    },
    {
        page: "/web/money/conversion/",
        text: "First, select source dust to convert by selecting the 'White dust' toggle. Then choose the value on the bar until the number right to it is equal to 2. This number represent an amount of dust to convert.",
    },
    {
        page: "/web/money/conversion/",
        text: "Finally, press the 'Convert' button to start the conversion.",
        trigger: "conversionStart"
    },
    {
        page: "/web/money/conversion/",
        text: "As the conversion has started, the selection bar is replaced with a progress bar. Wait for the conversion to end.",
        trigger: "conversionEnd"
    },
    {
        page: "/web/money/conversion/",
        text: "Now you can claim the result by clicking on the 'Claim' button.",
        trigger: "conversionClaim"
    },
    {
        page: "/web/money/conversion/",
        text: "Now let's talk about spending your dust. Click on 'Shop' link on the navigation bar.",
        nextPage: "/web/shop/"
    },
    {
        page: "/web/shop/",
        text: "Here you can buy chests with your dust. On purchase, each chest drops a random character of the matching rarity. It may drop a character you already have, so luck is a component of chests.",
    },
    {
        page: "/web/shop/",
        text: "This concludes the tutorial. Have a good time! Press 'Close' to close the tutorial.",
    },
]

var step = null;

function initTutorial() {
    const gotStep = window.localStorage.getItem("tutorial");
    if (!gotStep) {
        return;
    }

    step = parseInt(gotStep);
    showTutorialStep();
}

function stopTutorial() {
    step = null;
    window.localStorage.removeItem("tutorial");
    hideTutorial();
}

function showTutorialStep() {
    const scenarioStep = tutorialScenario[step];
    if (window.location.pathname !== scenarioStep.page) {
        if (window.location.pathname === scenarioStep.nextPage) {
            nextTutorialStep();
            return;
        }

        // Show return to tutorial prompt
        return;
    }

    const canNext = step === tutorialScenario.length - 1 || scenarioStep.trigger || scenarioStep.nextPage;
    showTutorialWindow(scenarioStep.text, canNext);
}

function showTutorialWindow(text, canNext) {
    const element = document.getElementById("tutorial");
    element.classList.remove("hidden");

    const textElement = document.getElementById("tutorial-text");
    textElement.innerText = text;

    document.getElementById("tutorial-next").disabled = canNext;
}

function hideTutorial() {
    const element = document.getElementById("tutorial");
    element.classList.add("hidden");
}

function enableTutorial() {
    window.localStorage.setItem("tutorial", "0");
}

function nextTutorialStep() {
    step++;
    window.localStorage.setItem("tutorial", step);

    if (step < tutorialScenario.length) {
        showTutorialStep();
    } else {
        stopTutorial();
    }
}

function tutorialTrigger(trigger) {
    if (!step) {
        return;
    }

    if (tutorialScenario[step].trigger === trigger) {
        nextTutorialStep();
    }
}

window.addEventListener("load", initTutorial);