import React from 'react'
import { FaTrain } from 'react-icons/fa'

const Navbar = () => {
    return (
        <div>
            {/* NAVBAR */}
            {/* NAVIGATION BAR */}
            <nav className="bg-white border-b border-gray-100 px-6 py-4">
                <div className="max-w-6xl mx-auto flex justify-between items-center">
                    <div className="flex items-center gap-2 font-black text-2xl tracking-tighter text-slate-800">
                        <div className="bg-slate-800 text-white p-1.5 rounded-lg">
                            <FaTrain size={20} />
                        </div>
                        RAIL<span className="text-orange-500">GO</span>
                    </div>
                    <div className="hidden md:flex gap-8 text-sm font-medium text-gray-600">
                        <a href="#" className="hover:text-slate-900 transition">Trains</a>
                        <a href="#" className="hover:text-slate-900 transition">PNR Status</a>
                        <a href="#" className="hover:text-slate-900 transition">Holidays</a>
                    </div>
                    <button className="text-sm font-semibold bg-gray-100 px-5 py-2 rounded-full hover:bg-gray-200 transition">
                        Login
                    </button>
                </div>
            </nav>
        </div>
    )
}

export default Navbar