'use client'
import React, { useEffect, useState } from 'react'
import { useSearchParams, useRouter } from 'next/navigation'
import { FaTrain, FaClock, FaArrowRight, FaFilter } from 'react-icons/fa'
import { apiFetch } from '@/lib/api'
import Navbar from '@/app/_components/Navbar'
import Footer from '@/app/_components/Footer'

interface Train {
    id: number
    trainnumber: number
    trainname: string
    source: string
    destination: string
    day: string
    arrivaltime: string
    departuretime: string
}

export default function SearchResultsPage() {
    const searchParams = useSearchParams()
    const from = searchParams.get('from')
    const to = searchParams.get('to')
    const date = searchParams.get('date')
    const [trains, setTrains] = useState<Train[]>([])
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState('')
    const router = useRouter()

    const handleBook = (train: Train, cls: string) => {
        const params = new URLSearchParams({
            train_id: train.id.toString(),
            train_name: train.trainname,
            date: date || '',
            class: cls
        })
        router.push(`/nget/book?${params.toString()}`)
    }

    useEffect(() => {
        const fetchTrains = async () => {
            try {
                // Since backend GetAllTrain is Admin only, we might hit an issue.
                // However, we'll try to fetch all and filter locally for now.
                const res = await apiFetch('/train/all-train')
                let filteredTrains = res.data

                if (from) {
                    filteredTrains = filteredTrains.filter((t: Train) =>
                        t.source.toLowerCase().includes(from.split(' ')[0].toLowerCase())
                    )
                }
                if (to) {
                    filteredTrains = filteredTrains.filter((t: Train) =>
                        t.destination.toLowerCase().includes(to.split(' ')[0].toLowerCase())
                    )
                }

                setTrains(filteredTrains)
            } catch (err: any) {
                setError(err.message)
            } finally {
                setLoading(false)
            }
        }

        fetchTrains()
    }, [from, to])

    return (
        <div className="min-h-screen bg-[#f7f7f5] flex flex-col">
            <Navbar />

            <main className="flex-1 max-w-6xl w-full mx-auto p-6 space-y-8">

                {/* Search Summary */}
                <div className="bg-white rounded-3xl p-8 border border-gray-100 shadow-[0_8px_30px_rgba(0,0,0,0.04)] flex flex-col md:flex-row items-center justify-between gap-6">
                    <div className="flex items-center gap-6">
                        <div className="text-left">
                            <div className="text-[10px] font-black uppercase tracking-[0.2em] text-gray-400 mb-1">From</div>
                            <div className="text-xl font-black text-gray-900">{from}</div>
                        </div>
                        <FaArrowRight className="text-gray-200" size={20} />
                        <div className="text-left">
                            <div className="text-[10px] font-black uppercase tracking-[0.2em] text-gray-400 mb-1">To</div>
                            <div className="text-xl font-black text-gray-900">{to}</div>
                        </div>
                    </div>

                    <div className="h-px w-full md:w-px md:h-12 bg-gray-100" />

                    <div className="flex items-center gap-12">
                        <div className="text-left">
                            <div className="text-[10px] font-black uppercase tracking-[0.2em] text-gray-400 mb-1">Departure</div>
                            <div className="text-base font-bold text-gray-600">{date}</div>
                        </div>
                        <button className="bg-gray-900 text-white px-6 py-3 rounded-2xl font-bold text-sm hover:bg-gray-800 transition-all">
                            Modify Search
                        </button>
                    </div>
                </div>

                {/* Results */}
                <div className="space-y-4">
                    <div className="flex items-center justify-between px-2">
                        <h3 className="text-lg font-black text-gray-900 mt-5">Available Trains ({trains.length})</h3>
                        <button className="flex items-center gap-2 text-gray-400 hover:text-gray-900 transition-colors font-bold text-sm">
                            <FaFilter size={12} />
                            Filters
                        </button>
                    </div>

                    {loading ? (
                        <div className="text-center py-20 text-gray-400 font-bold animate-pulse">Searching for best routes...</div>
                    ) : error ? (
                        <div className="bg-red-50 border border-red-100 text-red-500 p-6 rounded-3xl text-center font-bold">
                            {error === "Unauthorized" ? "Please sign in to search trains." : error}
                        </div>
                    ) : trains.length === 0 ? (
                        <div className="bg-white border border-gray-100 p-12 rounded-[32px] text-center space-y-4">
                            <div className="w-16 h-16 bg-gray-50 rounded-2xl flex items-center justify-center mx-auto mb-6">
                                <FaTrain className="text-gray-200" size={32} />
                            </div>
                            <p className="text-xl font-bold text-gray-900">No trains found</p>
                            <p className="text-gray-400 max-w-xs mx-auto text-sm leading-relaxed">We couldn't find any trains matching your route. Try searching for different stations or dates.</p>
                        </div>
                    ) : (
                        <div className="grid gap-4">
                            {trains.map((train) => (
                                <div key={train.id} className="bg-white rounded-[32px] border border-gray-100 p-8 shadow-[0_8px_30px_rgba(0,0,0,0.02)] hover:shadow-[0_8px_40px_rgba(0,0,0,0.06)] hover:border-gray-200 transition-all group">
                                    <div className="flex flex-col md:flex-row md:items-center justify-between gap-8">

                                        {/* Train Info */}
                                        <div className="space-y-4">
                                            <div className="flex items-center gap-3">
                                                <div className="bg-gray-900 text-white px-3 py-1 rounded-lg font-black text-[10px] tracking-wider uppercase">
                                                    {train.trainnumber}
                                                </div>
                                                <h4 className="text-xl font-black text-gray-900 group-hover:text-black transition-colors">{train.trainname}</h4>
                                            </div>
                                            <div className="flex items-center gap-8 text-gray-400">
                                                <div className="flex items-center gap-2">
                                                    <FaClock size={12} />
                                                    <span className="text-xs font-bold uppercase tracking-widest">{train.departuretime}</span>
                                                </div>
                                                <div className="w-px h-3 bg-gray-200 line-through" />
                                                <span className="text-xs font-bold uppercase tracking-widest">Runs on: {train.day}</span>
                                            </div>
                                        </div>

                                        {/* Classes/Pricing */}
                                        <div className="flex flex-wrap gap-3">
                                            {['SL', '3A', '2A', '1A'].map((cls) => (
                                                <button
                                                    key={cls}
                                                    className="px-6 py-4 rounded-2xl border border-gray-100 hover:border-gray-900 hover:bg-gray-900 hover:text-white transition-all text-left min-w-[120px] group/btn"
                                                >
                                                    <div className="text-[10px] font-black uppercase tracking-widest mb-1 group-hover/btn:text-gray-400">Class</div>
                                                    <div className="text-sm font-black">{cls}</div>
                                                    <div className="text-xs font-bold text-green-500 mt-2">Available</div>
                                                </button>
                                            ))}
                                        </div>

                                        {/* Action */}
                                        <button
                                            onClick={() => handleBook(train, 'SL')}
                                            className="bg-gray-900 text-white px-8 py-5 rounded-[20px] font-black text-sm hover:bg-gray-700 transition-all shadow-xl shadow-gray-200"
                                        >
                                            Book Now
                                        </button>
                                    </div>
                                </div>
                            ))}
                        </div>
                    )}
                </div>
            </main>

            <Footer />
        </div>
    )
}
