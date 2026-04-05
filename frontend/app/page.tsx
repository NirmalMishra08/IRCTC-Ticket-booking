'use client'
import React from 'react'
import { FaExchangeAlt } from 'react-icons/fa'
import { HiOutlineLocationMarker, HiOutlineCalendar, HiOutlineTicket } from 'react-icons/hi'
import Footer from './_components/Footer'
import { useRouter } from 'next/navigation'
import { useState } from 'react'

export default function Home() {
  const [from, setFrom] = useState('New Delhi (NDLS)')
  const [to, setTo] = useState('Mumbai Central (BCT)')
  const [date, setDate] = useState('2026-04-10')
  const [cls, setCls] = useState('All Classes')

  const router = useRouter()

  const swap = () => {
    setFrom(to)
    setTo(from)
  }

  const handleSearch = () => {
    const params = new URLSearchParams({
      from,
      to,
      date,
      class: cls
    })
    router.push(`/nget/search-results?${params.toString()}`)
  }

  return (
    <div className="min-h-screen flex flex-col bg-white">

      {/* ═══ HERO SECTION ═══ */}
      <section className="relative flex items-center justify-center min-h-[calc(100vh-68px)] px-6 overflow-hidden bg-[#f7f7f5]">

        {/* Background */}
        <div className="absolute inset-0 z-0">
          <div
            className="absolute inset-0 bg-cover bg-center opacity-[0.06]"
            style={{ backgroundImage: `url('/train-back.jpg')` }}
          />
          <div className="absolute inset-0 bg-gradient-to-b from-[#f7f7f5] via-[#f7f7f5]/80 to-[#f7f7f5]" />

          {/* Subtle radial glow behind the card */}
          <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[800px] h-[600px] bg-black/[0.02] rounded-full blur-[120px]" />
        </div>

        {/* Content */}
        <div className="relative z-10 max-w-3xl w-full text-center space-y-10">

          {/* Headline */}
          <div className="space-y-5 animate-fade-in-up">
            <h1 className="text-5xl md:text-7xl font-black tracking-tight leading-[1.05]">
              <span className="text-gray-900">Book Your</span>
              <br />
              <span className="text-gray-300">Next Journey</span>
            </h1>
            <p className="text-gray-400 text-base md:text-lg max-w-md mx-auto leading-relaxed">
              Search across the Indian Railways network.
              <br className="hidden md:block" />
              Fast. Secure. Reliable.
            </p>
          </div>

          {/* ═══ SEARCH CARD ═══ */}
          <div className="animate-fade-in-up delay-200">
            <div className="bg-white rounded-[28px] border border-gray-200 p-7 md:p-9 shadow-[0_8px_40px_rgba(0,0,0,0.08)] text-left">
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

                {/* Search Button */}
                <button
                  onClick={handleSearch}
                  className="w-full bg-gray-900 text-white py-4 rounded-2xl font-bold text-base hover:bg-gray-700 transition-all duration-300 shadow-[0_4px_20px_rgba(0,0,0,0.12)] hover:shadow-[0_4px_30px_rgba(0,0,0,0.18)] transform active:scale-[0.98]"
                >
                  Search Trains
                </button>
              </div>
            </div>
          </div>

          {/* Stats */}
          <div className="animate-fade-in delay-500 flex justify-center gap-10 text-center pt-2">
            <div>
              <div className="text-lg font-bold text-gray-600">14k+</div>
              <div className="text-xs text-gray-400">Daily Trains</div>
            </div>
            <div className="w-px bg-gray-200" />
            <div>
              <div className="text-lg font-bold text-gray-600">7.5k+</div>
              <div className="text-xs text-gray-400">Stations</div>
            </div>
            <div className="w-px bg-gray-200" />
            <div>
              <div className="text-lg font-bold text-gray-600">25M+</div>
              <div className="text-xs text-gray-400">Travelers</div>
            </div>
          </div>
        </div>
      </section>

      <Footer />
    </div>
  )
}