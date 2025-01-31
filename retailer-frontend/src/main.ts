import './style.css';

document.querySelector<HTMLDivElement>('#app')!.innerHTML = `
  <div class="container">
    <div class="glass-card">
      <div class="header">
        <h1 id="formTitle">Welcome Back</h1>
        <p id="formSubtitle">Please log in to your account</p>
      </div>
      <form id="authForm">
        <div class="input-group">
          <label for="username">Username</label>
          <input type="text" id="username" name="username" placeholder="Enter your username" required />
        </div>
        <div class="input-group">
          <label for="password">Password</label>
          <input type="password" id="password" name="password" placeholder="Enter your password" required />
        </div>
        <button type="submit" id="submitButton" class="login-button">Log In</button>
      </form>
      <p id="authMessage" class="message"></p>
      <div class="footer">
        <p id="toggleText">Don't have an account? <a href="#" id="toggleButton">Sign Up</a></p>
      </div>
    </div>
  </div>
`;

const authForm = document.querySelector<HTMLFormElement>('#authForm');
const authMessage = document.querySelector<HTMLParagraphElement>('#authMessage');
const toggleButton = document.querySelector<HTMLAnchorElement>('#toggleButton');
const formTitle = document.querySelector<HTMLHeadingElement>('#formTitle');
const formSubtitle = document.querySelector<HTMLParagraphElement>('#formSubtitle');
const submitButton = document.querySelector<HTMLButtonElement>('#submitButton');
const toggleText = document.querySelector<HTMLParagraphElement>('#toggleText');

let isLogin = true; // Toggle between login and registration

// Function to toggle between login and registration forms
const toggleForm = () => {
  isLogin = !isLogin;

  if (isLogin) {
    formTitle!.textContent = "Welcome Back";
    formSubtitle!.textContent = "Please log in to your account";
    submitButton!.textContent = "Log In";
    toggleText!.innerHTML = 'Don\'t have an account? <a href="#" id="toggleButton">Sign Up</a>';
  } else {
    formTitle!.textContent = "Create an Account";
    formSubtitle!.textContent = "Sign up to get started";
    submitButton!.textContent = "Sign Up";
    toggleText!.innerHTML = 'Already have an account? <a href="#" id="toggleButton">Log In</a>';
  }

  // Clear the form and message
  (document.querySelector<HTMLInputElement>('#username')!.value = ''),
  (document.querySelector<HTMLInputElement>('#password')!.value = ''),
  (authMessage!.textContent = '');
};

// Add event listener to toggle button
toggleButton?.addEventListener('click', (event) => {
  event.preventDefault();
  toggleForm();
});

// Add event listener to form submission
authForm?.addEventListener('submit', async (event) => {
  event.preventDefault();

  const username = (document.querySelector<HTMLInputElement>('#username')!.value).trim();
  const password = (document.querySelector<HTMLInputElement>('#password')!.value).trim();

  if (!username || !password) {
    authMessage!.textContent = "Please fill in all fields.";
    return;
  }

  const endpoint = isLogin ? "/login" : "/register";
  try {
    const response = await fetch(`https://retailer-x9lq.onrender.com${endpoint}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username, password }),
    });

    if (!response.ok) {
      throw new Error(isLogin ? 'Login failed' : 'Registration failed');
    }

    const data = await response.json();
    authMessage!.textContent = isLogin ? "Login successful!" : "Registration successful!";
    console.log(isLogin ? "JWT Token:" : "Response:", data);

    if (isLogin) {
      // Store the token in localStorage for login
      localStorage.setItem('token', data.token);
    } else {
      // Automatically switch to login after successful registration
      toggleForm();
    }
  } catch (error) {
    authMessage!.textContent = isLogin ? "Login failed. Please check your credentials." : "Registration failed. Please try again.";
    console.error(error);
  }
});