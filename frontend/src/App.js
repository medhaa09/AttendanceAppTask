import './App.css';
import {React,useEffect, useState} from 'react';
import LoginForm from './LoginForm';
import Register from './register';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import Home from './Home.js';
function App() {
  const [userRole,setUserRole]=useState('');
  useEffect(() => {
  const role=localStorage.getItem('role');
  setUserRole(role);
  },[]);
  return (
    <div className="App">
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<LoginForm />} />
          {userRole === 'admin' && (
            <Route path="/register" element={<Register />} />
          )}
           <Route path="/" element={<Home/>} />
        </Routes>
      </BrowserRouter>
    </div>
  );
}
export default App;
