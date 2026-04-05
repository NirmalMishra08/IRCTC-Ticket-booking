import React from 'react'
import { FaTrain } from 'react-icons/fa'

const Footer = () => {
    return (
        <footer className="mt-auto border-t border-[#1a1a1a] bg-[#050505]">
            <div className="max-w-6xl mx-auto px-6 py-10 flex flex-col md:flex-row justify-between items-center gap-6">
                {/* Brand */}
                <div className="flex items-center gap-2.5">
                    <div className="bg-white text-black p-1 rounded-md">
                        <FaTrain size={12} />
                    </div>
                    <span className="font-bold text-sm text-gray-500 tracking-tight">
                        RAIL<span className="text-gray-600">GO</span>
                    </span>
                    <span className="text-gray-700 text-xs ml-2">
                        © 2026 — An IRCTC Authorized Partner
                    </span>
                </div>

                {/* Links */}
                <div className="flex gap-6 text-xs font-medium text-gray-600">
                    <a href="#" className="hover:text-gray-300 transition-colors duration-300">
                        Privacy Policy
                    </a>
                    <a href="#" className="hover:text-gray-300 transition-colors duration-300">
                        Terms of Service
                    </a>
                    <a href="#" className="hover:text-gray-300 transition-colors duration-300">
                        Contact Us
                    </a>
                </div>
            </div>
        </footer>
    )
}

export default Footer