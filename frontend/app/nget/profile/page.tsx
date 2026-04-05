'use client'
import React, { useState } from 'react'
import { FaTrain, FaGoogle, FaEnvelope, FaLock, FaUser, FaArrowLeft } from 'react-icons/fa'
import { apiFetch } from '@/lib/api'
import { useRouter } from 'next/navigation'

export default function AuthPage() {
  const [isLogin, setIsLogin] = useState(true)
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [fullname, setFullname] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const router = useRouter()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      if (isLogin) {
        const res = await apiFetch('/auth/login', {
          method: 'POST',
          body: JSON.stringify({ email, password }),
        })
        localStorage.setItem('token', res.data.access_token)
        router.push('/')
      } else {
        // Backend for create-user was empty, but we'll try the common endpoint
        const res = await apiFetch('/auth/register', {
          method: 'POST',
          body: JSON.stringify({ email, password, fullname }),
        })
        setIsLogin(true)
        setError('Registration successful! Please login.')
      }
    } catch (err: any) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex h-screen w-full bg-white overflow-hidden">
      {/* Left Side Info */}
      <div className="w-1/2 hidden md:flex flex-col items-center justify-center bg-[#f7f7f5] p-12 relative overflow-hidden">
        <div className="absolute top-[-10%] right-[-10%] w-[400px] h-[400px] bg-black/[0.02] rounded-full blur-[80px]" />
        <div className="z-10 text-center space-y-6">
          <div className="w-20 h-20 bg-gray-900 rounded-3xl flex items-center justify-center mx-auto shadow-2xl mb-8 transform rotate-12">
            <FaTrain className="text-white text-4xl" />
          </div>
          <h1 className="text-4xl font-black text-gray-900 tracking-tight leading-tight">
            The Future of<br />Rail Travel.
          </h1>
          <p className="text-gray-400 max-w-xs mx-auto leading-relaxed">
            Experience seamless booking across India's vast railway network.
          </p>
        </div>
        <img
          className="absolute bottom-0 left-0 w-full opacity-10 pointer-events-none"
          src="https://raw.githubusercontent.com/prebuiltui/prebuiltui/main/assets/login/leftSideImage.png"
          alt="decoration"
        />
      </div>

      {/* Right Side Form */}
      <div className="w-full md:w-1/2 flex flex-col items-center justify-center p-8 bg-white">
        <form onSubmit={handleSubmit} className="w-full max-w-sm space-y-6 animate-fade-in-up">
          <div className="space-y-2 text-center md:text-left">
            <h2 className="text-3xl font-black text-gray-900 tracking-tight">
              {isLogin ? 'Sign In' : 'Create Account'}
            </h2>
            <p className="text-sm text-gray-400 font-medium">
              {isLogin
                ? 'Welcome back! Please sign in to continue'
                : 'Join us to start your journey across India'}
            </p>
          </div>

          {error && (
            <div className="bg-red-50 border border-red-100 text-red-500 px-4 py-3 rounded-2xl text-xs font-bold animate-shake">
              {error}
            </div>
          )}

          <div className="space-y-4">
            {!isLogin && (
              <div className="relative group">
                <div className="absolute left-5 top-1/2 -translate-y-1/2 text-gray-400 group-focus-within:text-gray-900 transition-colors">
                  <FaUser size={14} />
                </div>
                <input
                  type="text"
                  placeholder="Full Name"
                  className="w-full h-14 bg-gray-50 border border-gray-100 rounded-2xl pl-12 pr-6 outline-none focus:border-gray-900 focus:ring-1 focus:ring-gray-900 transition-all text-sm font-bold placeholder:text-gray-300"
                  required
                  value={fullname}
                  onChange={(e) => setFullname(e.target.value)}
                />
              </div>
            )}

            <div className="relative group">
              <div className="absolute left-5 top-1/2 -translate-y-1/2 text-gray-400 group-focus-within:text-gray-900 transition-colors">
                <FaEnvelope size={14} />
              </div>
              <input
                type="email"
                placeholder="Email Address"
                className="w-full h-14 bg-gray-50 border border-gray-100 rounded-2xl pl-12 pr-6 outline-none focus:border-gray-900 focus:ring-1 focus:ring-gray-900 transition-all text-sm font-bold placeholder:text-gray-300"
                required
                value={email}
                onChange={(e) => setEmail(e.target.value)}
              />
            </div>

            <div className="relative group">
              <div className="absolute left-5 top-1/2 -translate-y-1/2 text-gray-400 group-focus-within:text-gray-900 transition-colors">
                <FaLock size={14} />
              </div>
              <input
                type="password"
                placeholder="Password"
                className="w-full h-14 bg-gray-50 border border-gray-100 rounded-2xl pl-12 pr-6 outline-none focus:border-gray-900 focus:ring-1 focus:ring-gray-900 transition-all text-sm font-bold placeholder:text-gray-300"
                required
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
            </div>
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full h-14 rounded-2xl text-white bg-gray-900 hover:bg-gray-800 disabled:bg-gray-400 transition-all font-bold shadow-[0_4px_20px_rgba(0,0,0,0.12)] transform active:scale-[0.98] flex items-center justify-center"
          >
            {loading ? 'Processing...' : (isLogin ? 'Sign In' : 'Register')}
          </button>

          <div className="relative">
            <div className="absolute inset-0 flex items-center">
              <div className="w-full border-t border-gray-100"></div>
            </div>
            <div className="relative flex justify-center text-xs uppercase">
              <span className="bg-white px-4 text-gray-400 font-bold tracking-widest leading-none">OR</span>
            </div>
          </div>

          <button
            type="button"
            className="w-full h-14 bg-white border border-gray-100 rounded-2xl flex items-center justify-center gap-3 hover:bg-gray-50 transition-all font-bold text-sm text-gray-700 shadow-sm"
          >
            <FaGoogle className="text-gray-400" />
            Continue with Google
          </button>

          <p className="text-center text-xs font-bold text-gray-400 uppercase tracking-widest">
            {isLogin ? "Don't have an account?" : "Already have an account?"}
            <button
              type="button"
              onClick={() => setIsLogin(!isLogin)}
              className="ml-2 text-gray-900 hover:underline transition-colors"
            >
              {isLogin ? 'Register' : 'Sign In'}
            </button>
          </p>
        </form>
      </div>
    </div>
  )
}

