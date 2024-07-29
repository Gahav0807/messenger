'use client'
import "./styles.css"
import { toast,Toaster } from "sonner";
import { useState } from "react";
import { useRouter } from "next/navigation";

export default function AuthPage() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const router = useRouter();

  async function handleLogin() {
    try {
      const response = await fetch(`http://localhost:9000/tryAuth/${username}/${password}`);
  
      if (response.ok) {
        const data = await response.json();
        localStorage.setItem('messenger_token', data.token);      
        toast.success("Успешная авторизация!");
        
      } else {
        toast.error("Ошибка авторизации. Проверьте имя пользователя и пароль.");
      }
    } catch (error) {
      toast.error("Произошла ошибка при попытке авторизации.");
      console.error("Ошибка:", error);
    }
  }
  

  return (
    <div>
      <Toaster position="top-center" richColors />
      <div className='registartion-form'>
        <h1>Вход</h1>
        <p>Нет аккаунта? <a href='/register'>Регистрация</a></p>
        <div>
            <input
              type="text"
              placeholder="Имя"
              className='name-input'
              value={username}
              onChange={(e) => setUsername(e.target.value)}
            />
      
            <input
              type="password"
              placeholder="Пароль"
              className='password-input'
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
      
            <button onClick={handleLogin}>Next</button>
        </div>
      </div>
    </div>
  );
}
