'use client'
import React, { useEffect, useState } from 'react'
import { FaTrain, FaCalendarAlt, FaTicketAlt, FaClock, FaCheckCircle, FaExclamationCircle } from 'react-icons/fa'
import { apiFetch } from '@/lib/api'
import Navbar from '@/app/_components/Navbar'
import Footer from '@/app/_components/Footer'

interface Booking {
    id: number
    train_id: number
    train_name: string
    journey_date: string
    status: string
    booking_type: string
}

export default function MyBookingsPage() {
    const [bookings, setBookings] = useState<Booking[]>([])
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState('')

    useEffect(() => {
        const fetchBookings = async () => {
            try {
                // Backend has GetBookingbyUserId
                const res = await apiFetch('/booking/my-bookings') // Assuming this endpoint exists based on Querier
                setBookings(res.data || [])
            } catch (err: any) {
                setError(err.message)
            } finally {
                setLoading(false)
            }
        }
        fetchBookings()
    }, [])

    return (
        <div className="min-h-screen bg-[#f7f7f5] flex flex-col">
            <Navbar />

            <main className="flex-1 max-w-5xl w-full mx-auto p-6 py-12 space-y-10">
                <div className="space-y-2">
                    <h2 className="text-4xl font-black text-gray-900 tracking-tight">My Bookings</h2>
                    <p className="text-gray-400 font-bold uppercase tracking-widest text-xs">Manage your travel history and upcoming trips</p>
                </div>

                {loading ? (
                    <div className="py-20 text-center text-gray-400 font-bold animate-pulse">Retrieving your tickets...</div>
                ) : error ? (
                    <div className="bg-red-50 border border-red-100 text-red-500 p-8 rounded-[32px] text-center font-bold">
                        {error === "Unauthorized" ? "Please sign in to view your bookings." : error}
                    </div>
                ) : bookings.length === 0 ? (
                    <div className="bg-white border border-gray-100 p-16 rounded-[48px] text-center space-y-6 shadow-sm">
                        <div className="w-20 h-20 bg-gray-50 rounded-3xl flex items-center justify-center mx-auto mb-4">
                            <FaTicketAlt className="text-gray-200" size={32} />
                        </div>
                        <p className="text-xl font-black text-gray-900">No bookings yet</p>
                        <p className="text-gray-400 max-w-xs mx-auto text-sm font-medium leading-relaxed">You haven't booked any trains yet. Your future adventures will appear here!</p>
                        <button className="bg-gray-900 text-white px-8 py-4 rounded-2xl font-black text-sm hover:bg-gray-800 transition-all shadow-xl shadow-gray-100">
                            Book Your First Trip
                        </button>
                    </div>
                ) : (
                    <div className="grid gap-6">
                        {bookings.map((booking) => (
                            <div key={booking.id} className="bg-white rounded-[40px] border border-gray-100 p-8 shadow-[0_10px_40px_rgba(0,0,0,0.03)] hover:shadow-[0_10px_50px_rgba(0,0,0,0.06)] transition-all group">
                                <div className="flex flex-col md:flex-row gap-8 items-center justify-between">
                                    <div className="flex items-center gap-6">
                                        <div className="w-16 h-16 bg-gray-50 rounded-2xl flex items-center justify-center text-gray-900 group-hover:bg-gray-900 group-hover:text-white transition-all duration-500">
                                            <FaTrain size={24} />
                                        </div>
                                        <div className="space-y-1">
                                            <div className="flex items-center gap-3">
                                                <span className="text-xl font-black text-gray-900">{booking.train_name || `Train #${booking.train_id}`}</span>
                                                <div className="px-2 py-0.5 bg-gray-100 rounded text-[10px] font-black tracking-widest uppercase text-gray-400">
                                                    {booking.booking_type}
                                                </div>
                                            </div>
                                            <div className="flex items-center gap-4 text-gray-400 text-xs font-bold uppercase tracking-widest">
                                                <div className="flex items-center gap-1.5"><FaCalendarAlt size={10} /> {booking.journey_date}</div>
                                                <div className="w-px h-3 bg-gray-200" />
                                                <div className="flex items-center gap-1.5"><FaTicketAlt size={10} /> ID: #{booking.id}</div>
                                            </div>
                                        </div>
                                    </div>

                                    <div className="flex items-center gap-8">
                                        <div className="text-right">
                                            {booking.status === 'CONFIRMED' ? (
                                                <div className="flex items-center gap-2 text-green-500 font-black text-sm uppercase tracking-widest">
                                                    <FaCheckCircle />
                                                    Confirmed
                                                </div>
                                            ) : (
                                                <div className="flex items-center gap-2 text-orange-500 font-black text-sm uppercase tracking-widest">
                                                    <FaExclamationCircle />
                                                    {booking.status}
                                                </div>
                                            )}
                                        </div>
                                        <button className="bg-gray-50 text-gray-400 hover:bg-gray-900 hover:text-white px-6 py-4 rounded-2xl font-black text-xs uppercase tracking-widest transition-all">
                                            View Details
                                        </button>
                                    </div>
                                </div>
                            </div>
                        ))}
                    </div>
                )}
            </main>

            <Footer />
        </div>
    )
}
