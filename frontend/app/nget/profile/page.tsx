'use client'
import React, { useState } from 'react'
import { FaTrain, FaGoogle, FaEnvelope, FaLock, FaUser, FaArrowLeft } from 'react-icons/fa'

export default function AuthPage() {
  const [isLogin, setIsLogin] = useState(true)

  return (
    <div className="min-h-screen w-full flex items-center justify-center font-sans relative overflow-hidden bg-gray-50">
      
      {/* --- BACKGROUND LAYER (Matches your Hero) --- */}
      <div className="absolute inset-0 z-0">
        <div 
          className="absolute inset-0 bg-cover bg-center"
          style={{ backgroundImage: `url('/train-back.jpg')` }}
        />
        <div 
          className="absolute inset-0 backdrop-blur-md"
          style={{
            WebkitMaskImage: 'linear-gradient(to right, transparent 0%, black 40%, black 60%, transparent 100%)',
            maskImage: 'linear-gradient(to right, transparent 0%, black 40%, black 60%, transparent 100%)',
          }}
        />
        <div className="absolute inset-0 bg-white/60" />
      </div>

      {/* --- AUTH CARD --- */}
      <div className="relative z-10 w-full max-w-md px-6">
        
        {/* Back to Home Link */}
        <a href="/" className="inline-flex items-center gap-2 text-sm font-semibold text-slate-600 hover:text-slate-900 mb-8 transition-colors">
          <FaArrowLeft size={12} /> Back to Search
        </a>

        <div className="bg-white/80 backdrop-blur-xl rounded-[2.5rem] shadow-[0_20px_50px_rgba(0,0,0,0.1)] border border-white/20 p-10">
          
          {/* Logo */}
          <div className="flex justify-center mb-8">
            <div className="flex items-center gap-2 font-black text-2xl tracking-tighter text-slate-800">
              <div className="bg-slate-800 text-white p-1.5 rounded-lg">
                <FaTrain size={18} />
              </div>
              RAIL<span className="text-orange-500">GO</span>
            </div>
          </div>

          {/* Header */}
          <div className="text-center mb-8">
            <h2 className="text-2xl font-bold text-slate-900">
              {isLogin ? 'Welcome Back' : 'Create Account'}
            </h2>
            <p className="text-gray-500 text-sm mt-1">
              {isLogin ? 'Login to manage your bookings' : 'Join 25M+ travelers across India'}
            </p>
          </div>

          {/* Form */}
          <form className="space-y-4" onSubmit={(e) => e.preventDefault()}>
            {!isLogin && (
              <div className="relative">
                <FaUser className="absolute left-4 top-1/2 -translate-y-1/2 text-gray-400" size={14} />
                <input 
                  type="text" 
                  placeholder="Full Name"
                  className="w-full bg-white/50 border border-gray-200 rounded-2xl py-3.5 pl-11 pr-4 outline-none focus:border-slate-900 focus:ring-1 focus:ring-slate-900 transition-all font-medium text-slate-800"
                />
              </div>
            )}

            <div className="relative">
              <FaEnvelope className="absolute left-4 top-1/2 -translate-y-1/2 text-gray-400" size={14} />
              <input 
                type="email" 
                placeholder="Email Address"
                className="w-full bg-white/50 border border-gray-200 rounded-2xl py-3.5 pl-11 pr-4 outline-none focus:border-slate-900 focus:ring-1 focus:ring-slate-900 transition-all font-medium text-slate-800"
              />
            </div>

            <div className="relative">
              <FaLock className="absolute left-4 top-1/2 -translate-y-1/2 text-gray-400" size={14} />
              <input 
                type="password" 
                placeholder="Password"
                className="w-full bg-white/50 border border-gray-200 rounded-2xl py-3.5 pl-11 pr-4 outline-none focus:border-slate-900 focus:ring-1 focus:ring-slate-900 transition-all font-medium text-slate-800"
              />
            </div>

            {isLogin && (
              <div className="text-right">
                <button className="text-xs font-bold text-slate-500 hover:text-slate-900">Forgot Password?</button>
              </div>
            )}

            <button className="w-full bg-slate-900 text-white py-4 rounded-2xl font-bold hover:bg-slate-800 shadow-lg hover:shadow-slate-900/20 transition-all transform active:scale-[0.98] mt-2">
              {isLogin ? 'Sign In' : 'Get Started'}
            </button>
          </form>

          {/* Divider */}
          <div className="flex items-center my-6 gap-4">
            <div className="h-[1px] bg-gray-200 flex-1"></div>
            <span className="text-[10px] font-bold text-gray-400 uppercase tracking-widest">Or continue with</span>
            <div className="h-[1px] bg-gray-200 flex-1"></div>
          </div>

          {/* Social Auth */}
          <button className="w-full bg-white border border-gray-200 text-slate-700 py-3.5 rounded-2xl font-bold flex items-center justify-center gap-3 hover:bg-gray-50 transition-colors shadow-sm">
            <FaGoogle className="text-red-500" />
            Google
          </button>

          {/* Footer Toggle */}
          <p className="text-center mt-8 text-sm text-gray-500">
            {isLogin ? "Don't have an account?" : "Already have an account?"}{' '}
            <button 
              onClick={() => setIsLogin(!isLogin)}
              className="text-slate-900 font-bold hover:underline"
            >
              {isLogin ? 'Sign Up' : 'Log In'}
            </button>
          </p>
        </div>
      </div>
    </div>
  )
}