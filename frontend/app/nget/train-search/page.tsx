'use client'
import Footer from '@/app/_components/Footer'
import React, { useState } from 'react'
import { FaExchangeAlt, FaShieldAlt, FaClock, FaHeadset, FaBolt } from 'react-icons/fa'
import { HiOutlineLocationMarker, HiOutlineCalendar, HiOutlineTicket } from 'react-icons/hi'

export default function TrainSearchPage() {
  const [from, setFrom] = useState('New Delhi (NDLS)')
  const [to, setTo] = useState('Mumbai Central (BCT)')
  const [date, setDate] = useState('2026-04-10')
  const [cls, setCls] = useState('All Classes')
  const [tatkal, setTatkal] = useState(false)

  const swap = () => {
    setFrom(to)
    setTo(from)
  }

  const features = [
    { title: 'Instant Booking', icon: <FaBolt className="text-gray-500" />, desc: 'Confirm tickets in seconds with real-time availability.' },
    { title: 'Secure Payment', icon: <FaShieldAlt className="text-gray-500" />, desc: 'End-to-end encrypted payment gateways.' },
    { title: 'Fast Refund', icon: <FaClock className="text-gray-500" />, desc: '72-hour guaranteed refund processing.' },
    { title: '24/7 Support', icon: <FaHeadset className="text-gray-500" />, desc: 'Round the clock customer assistance.' },
  ]

  return (
    <div className="min-h-screen flex flex-col bg-white">

      {/* ═══ HERO ═══ */}
      <section className="relative pt-16 pb-32 px-6 overflow-hidden min-h-[620px] flex items-center bg-[#f7f7f5]">
        {/* Background */}
        <div className="absolute inset-0 z-0">
          <div
            className="absolute inset-0 bg-cover bg-center opacity-[0.06]"
            style={{ backgroundImage: `url('/train-back.jpg')` }}
          />
          <div className="absolute inset-0 bg-gradient-to-b from-[#f7f7f5] via-[#f7f7f5]/80 to-[#f7f7f5]" />
          <div className="absolute top-1/3 right-1/4 w-[600px] h-[400px] bg-black/[0.02] rounded-full blur-[100px]" />
        </div>

        <div className="max-w-6xl mx-auto relative z-10 flex flex-col lg:flex-row gap-12 items-center w-full">

          {/* Left Content */}
          <div className="flex-1 space-y-6 animate-fade-in-up">
            <div className="inline-flex items-center gap-2 bg-black/[0.04] border border-black/10 text-gray-500 px-4 py-1.5 rounded-full text-xs font-bold uppercase tracking-wider">
              <span className="relative flex h-2 w-2">
                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-black/30 opacity-75"></span>
                <span className="relative inline-flex rounded-full h-2 w-2 bg-black/70"></span>
              </span>
              Booking Live for Summer 2026
            </div>

            <h1 className="text-5xl lg:text-6xl font-black leading-[1.08] tracking-tight">
              <span className="text-gray-900">Your journey</span>
              <br />
              <span className="text-gray-300">starts with a click.</span>
            </h1>

            <p className="text-base text-gray-400 max-w-md leading-relaxed">
              Experience the most reliable way to book train tickets across the Indian Railways network.
            </p>

            <div className="flex gap-8 pt-5 border-t border-gray-200">
              {[
                { value: '14k+', label: 'Daily Trains' },
                { value: '7.5k+', label: 'Stations' },
                { value: '25M+', label: 'Users' },
              ].map((stat) => (
                <div key={stat.label}>
                  <div className="text-xl font-bold text-gray-700">{stat.value}</div>
                  <div className="text-xs text-gray-400 font-medium">{stat.label}</div>
                </div>
              ))}
            </div>
          </div>

          {/* Search Card */}
          <div className="w-full max-w-xl animate-fade-in-up delay-200">
            <div className="bg-white rounded-[28px] border border-gray-200 p-8 shadow-[0_8px_40px_rgba(0,0,0,0.08)]">
              <div className="space-y-5">

                {/* From / To */}
                <div className="relative flex flex-col md:flex-row gap-3">
                  <div className="flex-1 group border border-gray-200 rounded-2xl p-5 hover:border-gray-300 focus-within:border-gray-400 focus-within:ring-1 focus-within:ring-gray-200 transition-all duration-300 bg-gray-50/50">
                    <div className="flex items-center gap-2 text-gray-400 mb-1">
                      <HiOutlineLocationMarker size={14} />
                      <span className="text-[11px] font-bold uppercase tracking-[0.15em]">From</span>
                    </div>
                    <input
                      className="w-full outline-none font-bold text-gray-900 text-base bg-transparent placeholder:text-gray-300 mt-1"
                      value={from}
                      onChange={(e) => setFrom(e.target.value)}
                    />
                  </div>

                  <button
                    onClick={swap}
                    className="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 md:static md:translate-x-0 md:translate-y-0 w-10 h-10 bg-gray-100 border border-gray-200 rounded-full flex items-center justify-center hover:bg-gray-900 hover:text-white hover:border-gray-900 transition-all duration-300 z-10 shadow-sm text-gray-500"
                  >
                    <FaExchangeAlt size={11} className="md:rotate-0 rotate-90" />
                  </button>

                  <div className="flex-1 group border border-gray-200 rounded-2xl p-5 hover:border-gray-300 focus-within:border-gray-400 focus-within:ring-1 focus-within:ring-gray-200 transition-all duration-300 bg-gray-50/50">
                    <div className="flex items-center gap-2 text-gray-400 mb-1">
                      <HiOutlineLocationMarker size={14} />
                      <span className="text-[11px] font-bold uppercase tracking-[0.15em]">To</span>
                    </div>
                    <input
                      className="w-full outline-none font-bold text-gray-900 text-base bg-transparent placeholder:text-gray-300 mt-1"
                      value={to}
                      onChange={(e) => setTo(e.target.value)}
                    />
                  </div>
                </div>

                {/* Date & Class */}
                <div className="flex flex-col md:flex-row gap-3">
                  <div className="flex-1 border border-gray-200 rounded-2xl p-5 hover:border-gray-300 transition-all duration-300 bg-gray-50/50">
                    <div className="flex items-center gap-2 text-gray-400 mb-1">
                      <HiOutlineCalendar size={14} />
                      <span className="text-[11px] font-bold uppercase tracking-[0.15em]">Departure Date</span>
                    </div>
                    <input
                      type="date"
                      className="w-full outline-none font-bold text-gray-900 text-base bg-transparent mt-1"
                      value={date}
                      onChange={(e) => setDate(e.target.value)}
                    />
                  </div>

                  <div className="flex-1 border border-gray-200 rounded-2xl p-5 hover:border-gray-300 transition-all duration-300 bg-gray-50/50">
                    <div className="flex items-center gap-2 text-gray-400 mb-1">
                      <HiOutlineTicket size={14} />
                      <span className="text-[11px] font-bold uppercase tracking-[0.15em]">Travel Class</span>
                    </div>
                    <select
                      className="w-full outline-none font-bold text-gray-900 text-base bg-transparent appearance-none mt-1"
                      value={cls}
                      onChange={(e) => setCls(e.target.value)}
                    >
                      <option className="bg-white text-gray-900">All Classes</option>
                      <option className="bg-white text-gray-900">Sleeper (SL)</option>
                      <option className="bg-white text-gray-900">AC 3 Tier (3A)</option>
                      <option className="bg-white text-gray-900">AC 2 Tier (2A)</option>
                      <option className="bg-white text-gray-900">Exec. Chair Car (EC)</option>
                    </select>
                  </div>
                </div>

                {/* Tatkal Toggle */}
                <div className="flex flex-wrap gap-2 pt-1">
                  <button
                    onClick={() => setTatkal(!tatkal)}
                    className={`px-4 py-1.5 rounded-full text-xs font-semibold border transition-all duration-300 ${tatkal
                      ? 'bg-gray-900 border-gray-900 text-white shadow-[0_0_15px_rgba(0,0,0,0.1)]'
                      : 'bg-transparent border-gray-200 text-gray-400 hover:border-gray-400'
                      }`}
                  >
                    Tatkal
                  </button>
                </div>

                {/* CTA */}
                <button className="w-full bg-gray-900 text-white py-4 rounded-2xl font-bold text-base hover:bg-gray-700 transition-all duration-300 shadow-[0_4px_20px_rgba(0,0,0,0.12)] hover:shadow-[0_4px_30px_rgba(0,0,0,0.18)] transform active:scale-[0.98]">
                  Search Trains
                </button>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* ═══ FEATURES ═══ */}
      <div className="max-w-6xl mx-auto grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5 px-6 py-16 relative z-20">
        {features.map((card, idx) => (
          <div
            key={idx}
            className={`bg-white border border-gray-100 p-6 rounded-2xl hover:border-gray-200 hover:shadow-md hover:translate-y-[-4px] transition-all duration-500 group animate-fade-in-up delay-${(idx + 1) * 100}`}
          >
            <div className="w-11 h-11 bg-gray-50 rounded-xl flex items-center justify-center mb-4 text-lg group-hover:bg-gray-100 transition-colors duration-300">
              {card.icon}
            </div>
            <div className="font-bold text-gray-800 mb-1 text-sm">{card.title}</div>
            <div className="text-xs text-gray-400 leading-relaxed">{card.desc}</div>
          </div>
        ))}
      </div>

      <Footer />
    </div>
  )
}