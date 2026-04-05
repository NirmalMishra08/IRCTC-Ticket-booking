'use client'
import React, { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { FaTrain, FaUser, FaTicketAlt, FaShieldAlt, FaArrowRight } from 'react-icons/fa'
import { apiFetch } from '@/lib/api'
import Navbar from '@/app/_components/Navbar'
import Footer from '@/app/_components/Footer'

export default function ReviewPage() {
    const router = useRouter()
    const [bookingData, setBookingData] = useState<any>(null)
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState('')

    useEffect(() => {
        const data = sessionStorage.getItem('pendingBooking')
        if (!data) {
            router.push('/')
            return
        }
        setBookingData(JSON.parse(data))
    }, [router])

    const handlePayment = async () => {
        setLoading(true)
        setError('')
        try {
            // Create actual booking in backend
            const res = await apiFetch('/booking/create-booking', {
                method: 'POST',
                body: JSON.stringify({
                    train_id: Number(bookingData.trainId),
                    journey_date: bookingData.date,
                    seats: [bookingData.seat.seat_id],
                    passengers: [{
                        name: bookingData.passenger.name,
                        age: Number(bookingData.passenger.age),
                        gender: bookingData.passenger.gender
                    }],
                    booking_type: 'NORMAL'
                })
            })

            // Redirect to confirmation with booking ID
            router.push(`/nget/confirmation?booking_id=${res.data.id}`)
        } catch (err: any) {
            setError(err.message)
        } finally {
            setLoading(false)
        }
    }

    if (!bookingData) return null

    return (
        <div className="min-h-screen bg-[#f7f7f5] flex flex-col">
            <Navbar />

            <main className="flex-1 max-w-3xl w-full mx-auto p-6 space-y-10 py-12">
                <div className="text-center space-y-3">
                    <h2 className="text-4xl font-black text-gray-900 tracking-tight">Review Your Booking</h2>
                    <p className="text-gray-400 font-bold uppercase tracking-widest text-xs">One last step before you're ready to go</p>
                </div>

                {error && (
                    <div className="bg-red-50 border border-red-100 text-red-500 p-6 rounded-3xl font-bold flex items-center gap-3">
                        <div className="w-10 h-10 bg-red-100 rounded-xl flex items-center justify-center flex-shrink-0">⚠️</div>
                        {error}
                    </div>
                )}

                <div className="bg-white rounded-[40px] border border-gray-100 shadow-[0_20px_50px_rgba(0,0,0,0.04)] overflow-hidden">

                    {/* Header */}
                    <div className="bg-gray-900 p-10 text-white flex items-center justify-between">
                        <div className="space-y-1">
                            <div className="flex items-center gap-3">
                                <FaTrain className="text-white/40" />
                                <span className="text-2xl font-black">{bookingData.trainName}</span>
                            </div>
                            <p className="text-white/40 text-xs font-bold uppercase tracking-widest">{bookingData.date} • {bookingData.cls}</p>
                        </div>
                        <div className="text-right">
                            <div className="text-[10px] font-black uppercase tracking-[0.2em] text-white/40 mb-1">Status</div>
                            <div className="text-sm font-bold text-green-400">Available</div>
                        </div>
                    </div>

                    <div className="p-10 space-y-12">

                        {/* Passenger Details */}
                        <section className="space-y-6">
                            <div className="flex items-center gap-3">
                                <div className="w-8 h-8 bg-gray-50 rounded-lg flex items-center justify-center text-gray-400">
                                    <FaUser size={12} />
                                </div>
                                <h4 className="text-sm font-black uppercase tracking-widest text-gray-900">Passenger Info</h4>
                            </div>
                            <div className="bg-gray-50 rounded-3xl p-6 flex items-center justify-between">
                                <div>
                                    <div className="text-lg font-black text-gray-900">{bookingData.passenger.name}</div>
                                    <div className="text-xs font-bold text-gray-400 uppercase tracking-widest mt-1">{bookingData.passenger.age} years • {bookingData.passenger.gender}</div>
                                </div>
                                <div className="text-right">
                                    <div className="text-xs font-bold text-gray-400 uppercase tracking-widest mb-1">Seat</div>
                                    <div className="text-xl font-black text-gray-900">#{bookingData.seat.seat_no}</div>
                                    <div className="text-[10px] font-bold text-gray-400 uppercase">{bookingData.seat.berth}</div>
                                </div>
                            </div>
                        </section>

                        {/* Payment Summary */}
                        <section className="space-y-6">
                            <div className="flex items-center gap-3">
                                <div className="w-8 h-8 bg-gray-50 rounded-lg flex items-center justify-center text-gray-400">
                                    <FaTicketAlt size={12} />
                                </div>
                                <h4 className="text-sm font-black uppercase tracking-widest text-gray-900">Payment Details</h4>
                            </div>
                            <div className="space-y-4">
                                <div className="flex justify-between items-center text-sm font-bold text-gray-400">
                                    <span>Fare ({bookingData.cls})</span>
                                    <span className="text-gray-900 font-black">₹{bookingData.cls === 'SL' ? '450' : '1250'}</span>
                                </div>
                                <div className="flex justify-between items-center text-sm font-bold text-gray-400">
                                    <span>Booking Fee</span>
                                    <span className="text-gray-900 font-black">₹45</span>
                                </div>
                                <div className="pt-6 border-t border-gray-100 flex justify-between items-end">
                                    <div>
                                        <div className="text-[10px] font-black uppercase tracking-widest text-gray-400 mb-1">Total Amount</div>
                                        <div className="text-4xl font-black text-gray-900">₹{Number(bookingData.cls === 'SL' ? '450' : '1250') + 45}</div>
                                    </div>
                                    <button
                                        onClick={handlePayment}
                                        disabled={loading}
                                        className="bg-gray-900 text-white px-10 py-5 rounded-2xl font-black flex items-center gap-3 hover:bg-gray-800 transition-all shadow-xl shadow-gray-200 active:scale-95 disabled:bg-gray-400"
                                    >
                                        {loading ? 'Processing...' : 'Pay & Confirm'}
                                        <FaArrowRight />
                                    </button>
                                </div>
                            </div>
                        </section>
                    </div>

                    <div className="bg-gray-50 p-6 flex flex-col sm:flex-row items-center justify-center gap-6 border-t border-gray-100">
                        <div className="flex items-center gap-2 text-[10px] font-bold text-gray-400 uppercase tracking-widest">
                            <FaShieldAlt className="text-green-500" />
                            100% Secure Checkout
                        </div>
                    </div>
                </div>
            </main>

            <Footer />
        </div>
    )
}
