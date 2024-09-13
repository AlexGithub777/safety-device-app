// Function to submit form data for register
function register() {
    // Get the form data
    var form = document.getElementById("registerForm");
    registerForm.onsubmit = function () {
        if (
            this.elements["password"].value !=
            this.elements["confirm-password"].value
        ) {
            alert("Password not match");
            return false;
        }
        if (
            this.elements["username"].value.length < 6 ||
            this.elements["password"].value.length < 6
        ) {
            alert("Username and password must has at least 6 characters");
            return false;
        }
        if (
            this.elements["username"].value == this.elements["password"].value
        ) {
            alert("Username must different password");
            return false;
        }
    };
}
