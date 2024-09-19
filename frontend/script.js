var clicked = false;
const questionLabel = document.getElementById("questionlabel");
const options = document.querySelectorAll(".answers")
const evaluationspan = document.getElementById("evaluationspan");
const meaningButton = document.getElementById("get-meaning");
const registerContainer = document.getElementById("rc");
const loginContainer = document.getElementById("lc")
const form = document.getElementById("feedback-form")
let scorel = document.getElementById("score");
const evaluationspanColor = evaluationspan.backgroundColor
let score = 0;
let fdata;
let submitted = false;
registerContainer.style.display = "none";
loginContainer.style.display = "none";
function getData() {
    fetch("/getivquestion")
        .then((response) => response.json())
        .then((data) => {
            console.log(data);
            questionLabel.textContent = data.question;
            options[0].textContent = data.options[0];
            options[1].textContent = data.options[1];
            options[2].textContent = data.options[2];
            scorel.textContent = `score:${score}`;
            fdata = data;

        })
        .catch((error) => {
            console.error("Error:", error);
        });
}
getData();
function ev0() {
    options.forEach(button => button.disabled = true);
    if (options[0].textContent == fdata.answer) {
        evaluationspan.textContenft = "right answer";
        score++;
        evaluationEffect("green")
        scorel.textContent = `score:${score}`;
    } else {
        evaluationspan.textContent = "wrong answer";
        score--
        evaluationEffect("red")
        scorel.textContent = `score:${score}`;
    }
}
function ev1() {
    options.forEach(button => button.disabled = true);
    if (options[1].textContent == fdata.answer) {
        evaluationspan.textContent = "right answer";
        score++;
        evaluationEffect("green")
        scorel.textContent = `score:${score}`;
    } else {
        evaluationspan.textContent = "wrong answer";
        score--
        evaluationEffect("red")
        scorel.textContent = `score:${score}`;
    }
}
function ev2() {
    options.forEach(button => button.disabled = true);
    if (options[2].textContent == fdata.answer) {
        evaluationspan.textContent = "right answer";
        score++;
        evaluationEffect("green")
        scorel.textContent = `score:${score}`;
    } else {
        evaluationspan.textContent = "wrong answer";
        score--
        evaluationEffect("red")
        scorel.textContent = `score:${score}`;
    }
}
function evaluationEffect(color) {
    evaluationspan.style.backgroundColor = color
    setTimeout(function () {
        evaluationspan.style.backgroundColor = evaluationspanColor
    }, 1000)
}
function ChangeAnswerState() {
    options.forEach(button => button.disabled = false);
}
function submitFeedback() {
    submitted = !submitted
    document.getElementById('hidden_iframe').onload = function () {
        if (submitted) {
            document.getElementById('response').innerText = "Thank you for your feedback 🙏";
            submitted = false;
        }
    }
}
function showMeaning() {
    meaningButton.textContent = fdata.meaning
}
function hideMeaning() {
    meaningButton.textContent = "Get french meaning"
}
function openLogindiv() {
    if (loginContainer.style.display === "none") {
        registerContainer.style.display = "none"
        loginContainer.style.display = "block";
    } else {
        loginContainer.style.display = "none";
    }
}
function openRegisterdiv() {
    if (registerContainer.style.display === "none") {
        loginContainer.style.display = "none";
        registerContainer.style.display = "block";
    } else {
        registerContainer.style.display = "none";
    }
}
function loginRequest() {
    const formData = new URLSearchParams();
    formData.append('Name', document.getElementById("login-Name").value);
    formData.append('Password', document.getElementById("login-Password").value);
    formData.append('score', document.getElementById("score-input").value);

    // Send the POST request with URL-encoded form data
    fetch("/login", {
        method: "POST",
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded' // Set the content type to URL-encoded form data
        },
        body: formData // Convert URLSearchParams to a URL-encoded string
    })
        .then(response => {
            return response.text().then(text => ({
                status: response.status,
                text: text
            }));
        })
        .then(({ status, text }) => {
            if (status === 200) {
                console.log("Login successful");
                openLogindiv()
                alert("Login successfull\nscore submited")
            } else {
                console.log("Status code is not 200 :(");
            }
            console.log('Response text:', text);
        })
        .catch(error => {
            console.error('Error:', error);
            alert("Error logging in , check your credentials")
        });
}
document.getElementById("login-form").addEventListener('submit', function (event) {
    event.preventDefault()
    console.log("hello")
    loginRequest()
})
function registerRequest() {
    const formData = new URLSearchParams();
    formData.append('Name', document.getElementById("register-Name").value);
    formData.append('Password', document.getElementById("register-Password").value);
    formData.append('CPassword', document.getElementById("register-CPassword").value);

    fetch("/register", {
        method: "POST",
        body: formData
    }).then(response => {
        return response.text().then(text => ({
            status: response.status,
            text: text
        }));
    }).then(({ status, text }) => {
        if (status === 200) {
            console.log("register successful");
            alert("Register successfull")
            openRegisterdiv()
        } else {
            console.log("Status code is not 200 :(");
        }
        console.log('Response text:', text);
    }).catch(error => {
        console.error('Error:', error);
    });
}
document.getElementById("register-form").addEventListener("submit", function (event) {
    event.preventDefault();
    if (document.getElementById("register-Password").value === document.getElementById("register-CPassword").value) {
        console.log("form valid");
        registerRequest();
    } else {
        alert("Error: Password and Confirm Password should be matched");
    }
});
function continuouslyUpdateInputValue() {
    setInterval(() => {
        console.log(document.getElementById('score-input').value)
        document.getElementById('score-input').value = Number(score)
    }, 100); // Update every 100 milliseconds
}
continuouslyUpdateInputValue();
document.getElementById("leaderboard-btn").addEventListener("click", function () {
    window.location.href = "/leaderboard"
})
