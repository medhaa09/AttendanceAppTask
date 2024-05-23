import React, { useState } from 'react';
import './register.css';
import axios from 'axios';

export default function Register() {
  const [name, setName] = useState("");
  const [pass, setPass] = useState("");
  const [handle, setHandle] = useState("");
  const [image, setImage] = useState(null);
  const [preview, setPreview] = useState(null);

  const handleImageChange = (e) => {
    setImage(e.target.files[0]);
    setPreview(URL.createObjectURL(e.target.files[0]));
  };

  const handlechange = (c) => {
    c.preventDefault();
    
    const formData = new FormData();
    formData.append('name', name);
    formData.append('handle', handle);
    formData.append('pass', pass);
    formData.append('image', image);
    console.log("formdata ",formData);
    axios.post('http://localhost:8080/register',formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
        'Authorization': 'Bearer ' + localStorage.getItem('token')
      },
      withCredentials: true 
    })
    .then((res) => {
      alert('Registration Successful');
    })
    .catch((error) => {
      alert('Error occurred: ' + error.message);
    });
  };

  return (
    <div className="extradiv">
      <div className="register-form">
        <div className="header">
          <h1>Register</h1>
          <p>Enter information to create the student's account</p>
        </div>
        <form className="form" onSubmit={handlechange}>
          <div className="form-group">
            <label htmlFor="name">Full Name</label>
            <input id="name" placeholder="Enter full name" onChange={(e) => setName(e.target.value)} required />
          </div>
          <div className="form-group">
            <label htmlFor="userid">User id</label>
            <input id="userid" placeholder="Enter user id" onChange={(e) => setHandle(e.target.value)} required />
          </div>
          <div className="form-group">
            <label htmlFor="password">Password</label>
            <input id="password" type="password" placeholder="Enter password" onChange={(e) => setPass(e.target.value)} required />
          </div>
          <div className="form-group">
            <label htmlFor="image">Image</label>
            <input id="image" type="file" accept=".jpg" onChange={handleImageChange} required />
            {preview && (
              <img src={preview} alt="Preview" className="preview" />
            )}
          </div>
          <button type="submit" className="submit-btn">Register</button>
        </form>
      </div>
    </div>
  );
}
