
import React, { useState } from 'react';
import './LoginForm.css'; 
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
const LoginForm = () => {
  const [handle, setHandle] = useState('');
  const [pass, setPassword] = useState('');
  const navigate=useNavigate();
  const handleSubmit = (e) => {
    e.preventDefault();
    axios.post('http://localhost:8080/login',{handle, pass})
    .then((res) => {
      const responseData = res.data;
      console.log(res.data);
      const { token, refreshToken, role } = responseData;
    
      localStorage.setItem('token', token);
      localStorage.setItem('refreshToken', refreshToken)
      localStorage.setItem('role', role)
     alert('login Successful')
     if (role==="admin"){
        navigate("/register")
     }
     else{
      navigate("/")
     }

     
    })
    .catch((error) => {
      alert('error occured:' + error.message)
    });
  };
  return (
    <div className="login-container">
      <form className="login-form" onSubmit={handleSubmit}>
        <h2>Login</h2>
        <div className="input-group">
          <label htmlFor="Userid">User id</label>
          <input
            id="Userid"
            required
            value={handle}
            onChange={(e) => setHandle(e.target.value)}
          />
        </div>
        <div className="input-group">
          <label htmlFor="password">Password</label>
          <input
            type="password"
            id="password"
            required
            value={pass}
            onChange={(e) => setPassword(e.target.value)}
          />
        </div>
        <button type="submit" className="login-button">Log In</button>
      </form>
    </div>
  );
};
export default LoginForm;