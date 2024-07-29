'use client'
import "./styles.css"
import { toast,Toaster } from "sonner";
import { useState } from "react";
import { useRouter } from "next/navigation";

export default function RegPage() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const router = useRouter();

  async function handleRegister() {
    try {
      const response = await fetch(`http://localhost:9000/tryRegister/${username}/${password}`);
  
      if (response.ok) {
        const data = await response.json();
        localStorage.setItem('messenger_token', data.token);      
        toast.success("Регистрация прошла успешно!");
        
      } else {
        toast.error("Ошибка регистрации. Пользователь с таким ником уже существует.");
      }
    } catch (error) {
      toast.error("Произошла ошибка при попытке регистрации.");
      console.error("Ошибка:", error);
    }
  }

  return (
    <div>
      <Toaster position="top-center" richColors />
      <div className='registartion-form'>
        <h1>Регистрация</h1>
        <p>Уже есть аккаунт? <a href='/auth'>Войти</a></p>
        <div>
            <input
              type="text"
              placeholder="Имя"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
            />
            <input
              type="password"
              placeholder="Пароль"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
            <button onClick={handleRegister}>Next</button>
        </div>
      </div>
    </div>
  );
}
