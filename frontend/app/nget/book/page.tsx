'use client'
import React, { useState, useEffect } from 'react'
import { useSearchParams, useRouter } from 'next/navigation'
import { FaUser, FaCouch, FaArrowRight, FaShieldAlt } from 'react-icons/fa'
import { apiFetch } from '@/lib/api'
import Navbar from '@/app/_components/Navbar'
import Footer from '@/app/_components/Footer'

interface Seat {
    seat_id: number
    seat_no: number
    berth: string
    status: string
}

export default function BookingPage() {
    const searchParams = useSearchParams()
    const router = useRouter()
    const trainId = searchParams.get('train_id')
    const trainName = searchParams.get('train_name')
    const date = searchParams.get('date')
    const cls = searchParams.get('class')

    const [seats, setSeats] = useState<Seat[]>([])
    const [selectedSeat, setSelectedSeat] = useState<number | null>(null)
    const [passenger, setPassenger] = useState({ name: '', age: '', gender: 'Male' })
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState('')

    useEffect(() => {
        const fetchSeats = async () => {
            try {
                const res = await apiFetch('/train/get-all-seats', {
                    method: 'POST',
                    body: JSON.stringify({ train_id: Number(trainId), travel_date: date })
                })
                setSeats(res.data || [])
            } catch (err: any) {
                setError(err.message)
            } finally {
                setLoading(false)
            }
        }
        if (trainId) fetchSeats()
    }, [trainId, date])

    const handleProceed = () => {
        if (!selectedSeat) return alert('Please select a seat')
        if (!passenger.name || !passenger.age) return alert('Please fill passenger details')

        // Store in session or pass via state
        const bookingData = {
            trainId,
            trainName,
            date,
            cls,
            seat: seats.find(s => s.seat_id === selectedSeat),
            passenger
        }
        sessionStorage.setItem('pendingBooking', JSON.stringify(bookingData))
        router.push('/nget/review')
    }

    return (
        <div className="min-h-screen bg-[#f7f7f5] flex flex-col">
            <Navbar />

            <main className="flex-1 max-w-6xl w-full mx-auto p-6 grid grid-cols-1 lg:grid-cols-3 gap-8">

                {/* Left: Seat Selection */}
                <div className="lg:col-span-2 space-y-6">
                    <div className="bg-white rounded-[32px] p-8 border border-gray-100 shadow-sm">
                        <h3 className="text-2xl font-black text-gray-900 mb-2">Select Your Seat</h3>
                        <p className="text-gray-400 text-sm font-medium mb-8 uppercase tracking-widest">{trainName} • {cls} • {date}</p>

                        {loading ? (
                            <div className="py-20 text-center text-gray-400 font-bold">Loading seat map...</div>
                        ) : error ? (
                            <div className="bg-red-50 text-red-500 p-6 rounded-2xl font-bold">{error}</div>
                        ) : (
                            <div className="grid grid-cols-4 sm:grid-cols-6 md:grid-cols-8 gap-3 max-w-2xl mx-auto">
                                {seats.map((seat) => (
                                    <button
                                        key={seat.seat_id}
                                        disabled={seat.status !== 'AVAILABLE'}
                                        onClick={() => setSelectedSeat(seat.seat_id)}
                                        className={`h-14 rounded-xl border flex flex-col items-center justify-center transition-all ${selectedSeat === seat.seat_id
                                                ? 'bg-gray-900 border-gray-900 text-white shadow-xl transform scale-110'
                                                : seat.status === 'AVAILABLE'
                                                    ? 'bg-white border-gray-100 hover:border-gray-900 text-gray-600'
                                                    : 'bg-gray-50 border-gray-100 text-gray-200 cursor-not-allowed'
                                            }`}
                                    >
                                        <span className="text-[10px] font-black">{seat.seat_no}</span>
                                        <span className="text-[8px] font-bold uppercase tracking-tighter">{seat.berth}</span>
                                    </button>
                                ))}
                            </div>
                        )}

                        <div className="mt-8 flex justify-center gap-6 text-xs font-bold uppercase tracking-widest text-gray-400">
                            <div className="flex items-center gap-2"><div className="w-3 h-3 bg-white border border-gray-100 rounded-sm" /> Available</div>
                            <div className="flex items-center gap-2"><div className="w-3 h-3 bg-gray-900 rounded-sm" /> Selected</div>
                            <div className="flex items-center gap-2"><div className="w-3 h-3 bg-gray-50 border border-gray-100 rounded-sm" /> Occupied</div>
                        </div>
                    </div>

                    {/* Passenger Form */}
                    <div className="bg-white rounded-[32px] p-8 border border-gray-100 shadow-sm space-y-6">
                        <h3 className="text-2xl font-black text-gray-900 mb-6">Passenger Details</h3>
                        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                            <div className="md:col-span-2 relative group">
                                <div className="absolute left-5 top-1/2 -translate-y-1/2 text-gray-300 group-focus-within:text-gray-900 transition-colors">
                                    <FaUser size={14} />
                                </div>
                                <input
                                    type="text"
                                    placeholder="Full Name"
                                    className="w-full h-14 bg-gray-50 border border-gray-100 rounded-2xl pl-12 pr-6 outline-none focus:border-gray-900 focus:ring-1 focus:ring-gray-900 transition-all text-sm font-bold placeholder:text-gray-300"
                                    value={passenger.name}
                                    onChange={e => setPassenger({ ...passenger, name: e.target.value })}
                                />
                            </div>
                            <input
                                type="number"
                                placeholder="Age"
                                className="w-full h-14 bg-gray-50 border border-gray-100 rounded-2xl px-6 outline-none focus:border-gray-900 focus:ring-1 focus:ring-gray-900 transition-all text-sm font-bold placeholder:text-gray-300"
                                value={passenger.age}
                                onChange={e => setPassenger({ ...passenger, age: e.target.value })}
                            />
                            <select
                                className="w-full h-14 bg-gray-50 border border-gray-100 rounded-2xl px-6 outline-none focus:border-gray-900 focus:ring-1 focus:ring-gray-900 transition-all text-sm font-bold text-gray-600 appearance-none bg-[url('https://cdn0.iconfinder.com/data/icons/user-interface-150/24/Chevron_Down-512.png')] bg-[length:12px] bg-[right_1.5rem_center] bg-no-repeat"
                                value={passenger.gender}
                                onChange={e => setPassenger({ ...passenger, gender: e.target.value })}
                            >
                                <option>Male</option>
                                <option>Female</option>
                                <option>Other</option>
                            </select>
                        </div>
                    </div>
                </div>

                {/* Right: Summary & Action */}
                <div className="space-y-6">
                    <div className="bg-gray-900 rounded-[32px] p-8 text-white space-y-8 sticky top-6 shadow-2xl">
                        <h3 className="text-xl font-black">Booking Summary</h3>

                        <div className="space-y-6 border-b border-white/10 pb-8">
                            <div className="flex justify-between items-start">
                                <div>
                                    <div className="text-[10px] font-black uppercase tracking-[0.2em] text-white/40 mb-1">Train</div>
                                    <div className="text-sm font-bold">{trainName}</div>
                                </div>
                                <div className="text-right">
                                    <div className="text-[10px] font-black uppercase tracking-[0.2em] text-white/40 mb-1">Date</div>
                                    <div className="text-sm font-bold">{date}</div>
                                </div>
                            </div>

                            <div>
                                <div className="text-[10px] font-black uppercase tracking-[0.2em] text-white/40 mb-1">Selected Seat</div>
                                <div className="text-2xl font-black">{selectedSeat ? seats.find(s => s.seat_id === selectedSeat)?.seat_no : 'None'}</div>
                                <div className="text-[10px] font-bold text-white/40 uppercase">{selectedSeat ? seats.find(s => s.seat_id === selectedSeat)?.berth : '-'} Berth</div>
                            </div>
                        </div>

                        <div className="space-y-4">
                            <div className="flex justify-between items-center">
                                <span className="text-white/40 text-xs font-bold uppercase tracking-widest">Base Fare</span>
                                <span className="font-black">₹{cls === 'SL' ? '450' : cls === '3A' ? '1250' : '1850'}</span>
                            </div>
                            <div className="flex justify-between items-center">
                                <span className="text-white/40 text-xs font-bold uppercase tracking-widest">Taxes & Fees</span>
                                <span className="font-black">₹45</span>
                            </div>
                            <div className="pt-4 border-t border-white/10 flex justify-between items-center">
                                <span className="text-lg font-black">Total</span>
                                <span className="text-2xl font-black">₹{Number(cls === 'SL' ? '450' : cls === '3A' ? '1250' : '1850') + 45}</span>
                            </div>
                        </div>

                        <button
                            onClick={handleProceed}
                            className="w-full bg-white text-gray-900 h-14 rounded-2xl font-black flex items-center justify-center gap-3 hover:bg-gray-100 transition-all transform active:scale-95 group shadow-xl"
                        >
                            Continue to Review
                            <FaArrowRight className="group-hover:translate-x-1 transition-transform" />
                        </button>

                        <div className="flex items-center justify-center gap-2 text-[10px] font-bold text-white/30 uppercase tracking-widest">
                            <FaShieldAlt />
                            Safe & Secure Checkout
                        </div>
                    </div>
                </div>
            </main>

            <Footer />
        </div>
    )
}
