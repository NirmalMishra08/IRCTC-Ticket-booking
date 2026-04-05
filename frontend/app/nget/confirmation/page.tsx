'use client'
import React from 'react'
import { useSearchParams, useRouter } from 'next/navigation'
import { FaCheckCircle, FaTrain, FaArrowRight, FaTicketAlt } from 'react-icons/fa'
import Navbar from '@/app/_components/Navbar'
import Footer from '@/app/_components/Footer'

export default function ConfirmationPage() {
    const searchParams = useSearchParams()
    const bookingId = searchParams.get('booking_id')
    const router = useRouter()

    return (
        <div className="min-h-screen bg-[#f7f7f5] flex flex-col">
            <Navbar />

            <main className="flex-1 flex items-center justify-center p-6 py-20">
                <div className="max-w-xl w-full bg-white rounded-[48px] p-12 text-center shadow-[0_20px_70px_rgba(0,0,0,0.06)] border border-gray-100 space-y-10 animate-fade-in-up">

                    <div className="space-y-6">
                        <div className="w-24 h-24 bg-green-50 rounded-[32px] flex items-center justify-center mx-auto shadow-sm">
                            <FaCheckCircle className="text-green-500 text-5xl" />
                        </div>
                        <div className="space-y-2">
                            <h2 className="text-4xl font-black text-gray-900 tracking-tight">Booking Successful!</h2>
                            <p className="text-gray-400 font-bold uppercase tracking-widest text-xs">Your journey has been confirmed</p>
                        </div>
                    </div>

                    <div className="bg-gray-50 rounded-[32px] p-8 space-y-6">
                        <div className="flex justify-between items-center text-left">
                            <div>
                                <div className="text-[10px] font-black uppercase tracking-widest text-gray-400 mb-1">Booking ID</div>
                                <div className="text-2xl font-black text-gray-900">#{bookingId}</div>
                            </div>
                            <div className="text-right">
                                <div className="text-[10px] font-black uppercase tracking-widest text-gray-400 mb-1">PNR Status</div>
                                <div className="text-base font-black text-green-500 uppercase tracking-widest">Confirmed</div>
                            </div>
                        </div>

                        <div className="pt-6 border-t border-gray-200/50 flex items-center justify-center gap-10">
                            <button
                                onClick={() => router.push('/nget/my-bookings')}
                                className="flex items-center gap-2 text-xs font-black uppercase tracking-widest text-gray-900 hover:text-gray-600 transition-colors"
                            >
                                <FaTicketAlt />
                                View Tickets
                            </button>
                            <div className="w-px h-4 bg-gray-200" />
                            <button className="flex items-center gap-2 text-xs font-black uppercase tracking-widest text-gray-900 hover:text-gray-600 transition-colors">
                                <FaTrain />
                                Manage Journey
                            </button>
                        </div>
                    </div>

                    <button
                        onClick={() => router.push('/')}
                        className="w-full bg-gray-900 text-white h-16 rounded-2xl font-black flex items-center justify-center gap-3 hover:bg-gray-800 transition-all shadow-xl shadow-gray-100 group"
                    >
                        Back to Home
                        <FaArrowRight className="group-hover:translate-x-1 transition-transform" />
                    </button>

                    <p className="text-gray-400 text-[10px] font-bold uppercase tracking-[0.2em] leading-relaxed">
                        A confirmation email has been sent to your registered address.<br />
                        Have a pleasant journey!
                    </p>
                </div>
            </main>

            <Footer />
        </div>
    )
}
