// Import the functions you need from the SDKs you need
import { initializeApp } from "firebase/app";
import { getAnalytics } from "firebase/analytics";
// TODO: Add SDKs for Firebase products that you want to use
// https://firebase.google.com/docs/web/setup#available-libraries

// Your web app's Firebase configuration
// For Firebase JS SDK v7.20.0 and later, measurementId is optional
const firebaseConfig = {
  apiKey: "AIzaSyC0y0Vt_mndRS9tLykL4dsoUUk5cGjdvwo",
  authDomain: "better-uptime-auth.firebaseapp.com",
  projectId: "better-uptime-auth",
  storageBucket: "better-uptime-auth.firebasestorage.app",
  messagingSenderId: "770634605145",
  appId: "1:770634605145:web:a2060e579763b9ad214748",
  measurementId: "G-8MGV4H3ZZN"
};

// Initialize Firebase
const app = initializeApp(firebaseConfig);
const analytics = getAnalytics(app);