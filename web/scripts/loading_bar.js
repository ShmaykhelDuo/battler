function LoadingBar(x, y, w, h, radius, id, c, c2) {
    this.transitionFrames = 30;
    this.transitionFramesLeft = 0;
    this.id = id;
    this.clickable = false;
    this.hoverable = false;
    this.rectColour = c;
    this.rightColour = c2;
    this.stopColour = c;
    this.percentage = 0.0;
    this.newPercentage = 0.0;
    this.x = x;
    this.y = y;
    this.width = w;
    this.w = 0;
    this.h = h;
    this.radius = radius;

    this.display = function () {
        if (this.transitionFramesLeft > 0) {
            this.percentage += (this.newPercentage - this.percentage) / this.transitionFramesLeft;
            this.transitionFramesLeft--;
        }
        this.w = (this.width * this.percentage / 100);
        this.stopColour = lerpColor(this.rectColour, this.rightColour, this.percentage / 100);
        this.setGradient(this.rectColour, this.stopColour);
        stroke(this.rectColour);
        strokeWeight(2);
        noFill();
        rect(this.x, this.y, this.width, this.h, this.radius);
    };

    this.setPercentage = function (perc) {
        if (this.transitionFramesLeft > 0) {
            this.newPercentage = perc;
        } else {
            this.percentage = perc;
        }
    };

    this.setNewPercentage = function (new_perc, fastUpdate) {
        this.newPercentage = new_perc;
        if (fastUpdate) {
            this.transitionFramesLeft = this.transitionFrames;
        }
    };

    this.setGradient = function (c1, c2) {
        let x = this.x;
        let y = this.y;
        let w = this.w;
        let h = this.h;
        let r = this.radius;
        noFill();
        strokeWeight(2);
        //circle at the beginning
        for (let i = x; (i < x + r) && (i < x + w); i++) {
            let inter = map(i, x, x + w, 0, 1);
            let c = lerpColor(c1, c2, inter);
            stroke(c);
            let top_y = y + r - sqrt(sq(r) - sq(i - x - r));
            let bot_y = y + h - r + sqrt(sq(r) - sq(i - x - r));
            line(i, top_y, i, bot_y);
        }

        for (let i = x + r; i < x + w - r; i++) {
            let inter = map(i, x, x + w, 0, 1);
            let c = lerpColor(c1, c2, inter);
            stroke(c);
            line(i, y, i, y + h);
        }
        //circle at the end
        for (let i = x + w; (i >= x + r) && (i >= x + w - r); i--) {
            let inter = map(i, x, x + w, 0, 1);
            let c = lerpColor(c1, c2, inter);
            stroke(c);
            let top_y = y + r - sqrt(sq(r) - sq(i - x - w + r));
            let bot_y = y + h - r + sqrt(sq(r) - sq(i - x - w + r));
            line(i, top_y, i, bot_y);
        }
    };

    this.in = function() {
        let x = mouseX;
        let y = mouseY;
        return (this.x <= x && x <= (this.width + this.x) && this.y <= y && y <= (this.h + this.y))
    };

    this.setColours = function (c1, c2) {
        this.rectColour = c1;
        this.rightColour = c2;
    };

    this.makeDraggable = function() {
        this.clickable = true;
        this.draggable = true;
    };

    this.makeNotDraggable = function() {
        this.clickable = false;
        this.draggable = false;
    };
}