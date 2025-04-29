import { showMessage, handlePostFetch, start} from "./main.js";

/*------ Authentication functions ------*/
export function submitSignUp(event){
    event.preventDefault();
    const form = document.getElementById('signUpForm');
    const formData = new FormData(form);

    const userFirstName = formData.get('firstName');
    const userLastName = formData.get('lastName');
    const userNickname = formData.get('nickName');
    const gender = formData.get('gender');
    const age = parseInt(formData.get('age'));
    const email = formData.get('email');
    const password = formData.get('password');
    const confirmPassword = formData.get('confirm-password');

    // Validate username
    if (!isValidName("First Name", userFirstName)) return;
    if (!isValidName("Last Name", userLastName)) return;
    if (!isValidName("Nickname", userNickname)) return;

    // Validate gender (check if one of the radio buttons is selected and its value is valid)
    if (!gender || !['Male', 'Female', 'Other'].includes(gender)) {
        return showMessage("Please select a valid gender (Male, Female, or Other).");
    }

    // validate age
    if (isNaN(age) || age < 9 || age > 100) {
        return showMessage("Please enter a valid age between 9 and 100.");
    }

    // Validate email
    if (email.trim() === "") {
        return showMessage("Email is required.");
    } else if (!/^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/.test(email)) {
        return showMessage("Please enter a valid email address.");
    }

    // Validate password
    if (password.trim() === "") {
        return showMessage("Password is required.");
    } else if (password.length <= 7) {
        return showMessage("Password must be at least 8 characters long.");
    } else if (password.toLowerCase() === "password") {
        return showMessage("Password cannot be 'password'.");
    }
    if (confirmPassword.trim() === "" || password !== confirmPassword){
        return showMessage("Passwords do not match.");
    }

    handlePostFetch('/signup', {
        firstName: userFirstName,
        lastName: userLastName,
        nickName: userNickname,
        gender: gender,
        age: age,
        email: email,
        password: password,
        confirmPassword: confirmPassword,
    }, "Signup successful!", start);
};

const nickNameRadio = document.getElementById("nickNameField");
const nickNameGroup = document.getElementById("nickNameGroup");
const emailRadio = document.getElementById("emailField");
const emailGroup = document.getElementById("emailGroup");
nickNameRadio.addEventListener("change", toggleFields);
emailRadio.addEventListener("change", toggleFields);

function toggleFields() {
if (nickNameRadio.checked) {
    nickNameGroup.style.display = "flex";
    emailGroup.style.display = "none";
} else if (emailRadio.checked) {
    nickNameGroup.style.display = "none";
    emailGroup.style.display = "flex";
}
}

export function submitLogIn(event){
    event.preventDefault();
    const form = document.getElementById('logInForm');
    const formData = new FormData(form);

    let userNickname = formData.get('nickName');
    let email = formData.get('email');
    const password = formData.get('password');

    // Validate username
    if (nickNameRadio.checked){
        console.log(userNickname)
        if (!isValidName("Nickname", userNickname)) return;
    } else {userNickname = ""}

    // Validate email
    if (emailRadio.checked){
        if (email.trim() === "") {
        return showMessage("Email is required.");
        } else if (!/^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/.test(email)) {
        return showMessage("Please enter a valid email address.");
        }
    } else {email = ""}

    // Validate password
    if (password.trim() === "") {
        return showMessage("Password is required.");
    } else if (password.length <= 7) {
        return showMessage("Password must be at least 8 characters long.");
    } else if (password.toLowerCase() === "password") {
        return showMessage("Password cannot be 'password'.");
    }

    handlePostFetch('/login', {
        nickName: userNickname,
        email: email,
        password: password,
    }, "Log in successful!", start);
};

export function submitLogOut(event){
    event.preventDefault();
    handlePostFetch('/logout', {}, "Log out successful!", start);
};

function isValidName(field, data){
    if (data.trim() === ""){
        showMessage(`${field} is required.`);
        return false;
    } else if (data.length > 16 || data.length < 3 || !/^[\u0000-\u007F]+$/.test(data)){
        showMessage(`${field}  must be between 3 to 16 alphanumeric characters, '_' or '-'`);
        return false;
    }
    return true;
};
