'use client'
import Footer from '@/app/_components/Footer'
import React, { useState } from 'react'
import { FaTrain, FaExchangeAlt, FaShieldAlt, FaClock, FaHeadset, FaBolt } from 'react-icons/fa'
import { HiOutlineLocationMarker, HiOutlineCalendar, HiOutlineTicket } from 'react-icons/hi'

export default function IRCTCPage() {
  const [from, setFrom] = useState('New Delhi (NDLS)')
  const [to, setTo] = useState('Mumbai Central (BCT)')
  const [date, setDate] = useState('2026-03-26')
  const [cls, setCls] = useState('All Classes')
  const [tatkal, setTatkal] = useState(false)
  const [flexible, setFlexible] = useState(false)
  const [disability, setDisability] = useState(false)

  const swap = () => {
    setFrom(to)
    setTo(from)
  }

  const infoCards = [
    { title: 'Instant Booking', icon: <FaBolt className="text-orange-500" />, desc: 'Confirm tickets in seconds.' },
    { title: 'Secure Payment', icon: <FaShieldAlt className="text-emerald-500" />, desc: 'Fully encrypted gateways.' },
    { title: 'Fast Refund', icon: <FaClock className="text-blue-500" />, desc: '72-hour refund processing.' },
    { title: '24/7 Support', icon: <FaHeadset className="text-purple-500" />, desc: 'Always here to help you.' },
  ]

  return (
    <div className="min-h-screen bg-white-50 font-sans flex flex-col">



      {/* HERO SECTION */}
      <section className="relative my-10 pt-16 pb-32 px-6 overflow-hidden min-h-[600px] flex items-center">
        <div className="absolute inset-0 z-0">
          {/* 1. Base Image (Sharp) */}
          <div
            className="absolute inset-0 bg-cover bg-center"
            style={{
              backgroundImage: `url('/train-back.jpg')`,
            }}
          />

          {/* 2. Blurry Layer (Middle Only) */}
          <div
            className="absolute inset-0 bg-cover bg-center backdrop-blur-md"
            style={{
              backgroundImage: `url('/train-back.jpg')`,
              // Mask: Transparent at edges (0-20%), Solid Black in middle (40-60%), Transparent at end (80-100%)
              WebkitMaskImage: 'linear-gradient(to right, transparent 0%, black 40%, black 60%, transparent 100%)',
              maskImage: 'linear-gradient(to right, transparent 0%, black 40%, black 60%, transparent 100%)',
            }}
          />

          {/* 3. Slight Darkening Overlay (Keeps text readable) */}
          <div className="absolute inset-0 bg-white/60 z-[1]" />
        </div>
        <div className="max-w-6xl mx-auto flex flex-col lg:flex-row gap-12 items-center">

          {/* LEFT CONTENT */}
          <div className="flex-1 space-y-6">
            <div className="inline-flex items-center gap-2 bg-orange-50 text-orange-700 px-4 py-1.5 rounded-full text-xs font-bold uppercase tracking-wider">
              <span className="relative flex h-2 w-2">
                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-orange-400 opacity-75"></span>
                <span className="relative inline-flex rounded-full h-2 w-2 bg-orange-500"></span>
              </span>
              Booking Live for Summer 2026
            </div>

            <h1 className="text-5xl lg:text-6xl font-extrabold text-slate-900 leading-[1.1]">
              Your journey <br />
              <span className="text-gray-400">starts with a click.</span>
            </h1>

            <p className="text-lg text-gray-500 max-w-md leading-relaxed">
              Experience the most reliable way to book train tickets across the Indian Railways network.
            </p>

            <div className="flex gap-8 pt-4 border-t border-gray-100">
              <div>
                <div className="text-2xl font-bold text-slate-800">14k+</div>
                <div className="text-sm text-gray-400 font-medium">Daily Trains</div>
              </div>
              <div>
                <div className="text-2xl font-bold text-slate-800">7.5k+</div>
                <div className="text-sm text-gray-400 font-medium">Stations</div>
              </div>
              <div>
                <div className="text-2xl font-bold text-slate-800">25M+</div>
                <div className="text-sm text-gray-400 font-medium">Users</div>
              </div>
            </div>
          </div>

          {/* SEARCH CARD */}
          <div className="w-full max-w-xl bg-white rounded-3xl shadow-[0_20px_50px_rgba(0,0,0,0.1)] border border-gray-100 p-8 relative z-10">
            <div className="space-y-5">

              {/* Station Selection */}
              <div className="relative flex flex-col md:flex-row gap-2">
                <div className="flex-1 group transition-all border border-gray-200 rounded-2xl p-4 hover:border-slate-400 focus-within:border-slate-900 focus-within:ring-1 focus-within:ring-slate-900">
                  <div className="flex items-center gap-2 text-gray-400 mb-1">
                    <HiOutlineLocationMarker size={16} />
                    <span className="text-[10px] font-bold uppercase tracking-widest">From</span>
                  </div>
                  <input
                    className="w-full outline-none font-bold text-slate-800 bg-transparent"
                    value={from}
                    onChange={(e) => setFrom(e.target.value)}
                  />
                </div>

                <button
                  onClick={swap}
                  className="absolute  left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 md:static md:translate-x-0 md:translate-y-0 w-10 h-10 bg-white border border-gray-200 rounded-full flex text-center items-center justify-center hover:bg-slate-900 hover:text-white transition-all shadow-md z-10"
                >
                  <FaExchangeAlt size={12} className="md:rotate-0 rotate-90" />
                </button>

                <div className="flex-1 group transition-all border border-gray-200 rounded-2xl p-4 hover:border-slate-400 focus-within:border-slate-900 focus-within:ring-1 focus-within:ring-slate-900">
                  <div className="flex items-center gap-2 text-gray-400 mb-1">
                    <HiOutlineLocationMarker size={16} />
                    <span className="text-[10px] font-bold uppercase tracking-widest">To</span>
                  </div>
                  <input
                    className="w-full outline-none font-bold text-slate-800 bg-transparent"
                    value={to}
                    onChange={(e) => setTo(e.target.value)}
                  />
                </div>
              </div>

              {/* Date & Class */}
              <div className="flex flex-col md:flex-row gap-4">
                <div className="flex-1 border border-gray-200 rounded-2xl p-4 hover:border-slate-400 transition-all">
                  <div className="flex items-center gap-2 text-gray-400 mb-1">
                    <HiOutlineCalendar size={16} />
                    <span className="text-[10px] font-bold uppercase tracking-widest">Departure Date</span>
                  </div>
                  <input
                    type="date"
                    className="w-full outline-none font-bold text-slate-800 bg-transparent"
                    value={date}
                    onChange={(e) => setDate(e.target.value)}
                  />
                </div>

                <div className="flex-1 border border-gray-200 rounded-2xl p-4 hover:border-slate-400 transition-all">
                  <div className="flex items-center gap-2 text-gray-400 mb-1">
                    <HiOutlineTicket size={16} />
                    <span className="text-[10px] font-bold uppercase tracking-widest">Travel Class</span>
                  </div>
                  <select
                    className="w-full outline-none font-bold text-slate-800 bg-transparent appearance-none"
                    value={cls}
                    onChange={(e) => setCls(e.target.value)}
                  >
                    <option>All Classes</option>
                    <option>Sleeper (SL)</option>
                    <option>AC 3 Tier (3A)</option>
                    <option>AC 2 Tier (2A)</option>
                    <option>Exec. Chair Car (EC)</option>
                  </select>
                </div>
              </div>

              {/* Toggles */}
              <div className="flex flex-wrap gap-2 pt-2">
                {[
                  { id: 'tatkal', label: 'Tatkal', state: tatkal, setter: setTatkal }
                ].map((tag) => (
                  <button
                    key={tag.id}
                    onClick={() => tag.setter(!tag.state)}
                    className={`px-4 py-1.5 rounded-full text-xs font-semibold border transition-all ${tag.state
                      ? 'bg-slate-900 border-slate-900 text-white shadow-md'
                      : 'bg-white border-gray-200 text-gray-500 hover:border-gray-400'
                      }`}
                  >
                    {tag.label}
                  </button>
                ))}
              </div>

              <button className="w-full bg-slate-900 text-white py-4 rounded-2xl font-bold text-lg hover:bg-slate-800 transition-all shadow-lg hover:shadow-xl transform active:scale-[0.98]">
                Search Trains
              </button>
            </div>
          </div>
        </div>
      </section>

      {/* FEATURE CARDS */}
      <div className="max-w-6xl mx-auto grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 px-6  py-6 relative z-20">
        {infoCards.map((card, idx) => (
          <div key={idx} className="bg-white p-6 rounded-2xl shadow-xl border border-gray-50 hover:translate-y-[-5px] transition-transform duration-300">
            <div className="w-12 h-12 bg-gray-50 rounded-xl flex items-center justify-center mb-4 text-xl">
              {card.icon}
            </div>
            <div className="font-bold text-slate-800 mb-1">{card.title}</div>
            <div className="text-sm text-gray-500">{card.desc}</div>
          </div>
        ))}
      </div>

      <Footer />


    </div>
  )
}