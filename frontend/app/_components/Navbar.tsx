'use client'
import React, { useState, useEffect } from 'react'
import { FaTrain } from 'react-icons/fa'
import { HiOutlineMenuAlt3, HiX } from 'react-icons/hi'

const Navbar = () => {
    const [scrolled, setScrolled] = useState(false)
    const [mobileOpen, setMobileOpen] = useState(false)

    useEffect(() => {
        const handleScroll = () => setScrolled(window.scrollY > 20)
        window.addEventListener('scroll', handleScroll)
        return () => window.removeEventListener('scroll', handleScroll)
    }, [])

    const navLinks = [
        { label: 'Trains', href: '/' },
        { label: 'PNR Status', href: '#' },
        { label: 'My Bookings', href: '/nget/my-bookings' },
    ]

    return (
        <>
            <nav
                className={`fixed top-0 left-0 right-0 z-50 transition-all duration-500 ${scrolled
                    ? 'glass shadow-lg'
                    : 'bg-[#0a0a0a]'
                    }`}
            >
                <div className="max-w-6xl mx-auto px-6 py-4 flex justify-between items-center">
                    {/* Logo */}
                    <a href="/" className="flex items-center gap-2.5 group">
                        <div className="bg-white text-black p-1.5 rounded-lg transition-transform duration-300 group-hover:scale-110">
                            <FaTrain size={16} />
                        </div>
                        <span className={`font-black text-xl tracking-tight text-white`}>
                            RAIL<span className="text-gray-500">GO</span>
                        </span>
                    </a>

                    {/* Desktop Nav */}
                    <div className="hidden md:flex items-center gap-8">
                        {navLinks.map((link) => (
                            <a
                                key={link.label}
                                href={link.href}
                                className="text-sm font-medium text-gray-400 hover:text-white transition-colors duration-300 relative group"
                            >
                                {link.label}
                                <span className="absolute -bottom-1 left-0 w-0 h-[1px]  transition-all duration-300 group-hover:w-full" />
                            </a>
                        ))}
                    </div>

                    {/* Desktop CTA */}
                    <div className="hidden md:flex items-center gap-3">
                        <a
                            href="/nget/profile"
                            className="text-sm font-semibold px-5 py-2 rounded-full border border-gray-700 text-gray-300 hover:border-white hover:text-white transition-all duration-300 hover:shadow-[0_0_20px_rgba(255,255,255,0.08)]"
                        >
                            Login
                        </a>
                    </div>

                    {/* Mobile Menu Toggle */}
                    <button
                        onClick={() => setMobileOpen(!mobileOpen)}
                        className="md:hidden text-gray-400 hover:text-white transition-colors"
                    >
                        {mobileOpen ? <HiX size={24} /> : <HiOutlineMenuAlt3 size={24} />}
                    </button>
                </div>

                {/* Gradient bottom line */}
                <div className={`gradient-line transition-opacity duration-500 ${scrolled ? 'opacity-100' : 'opacity-0'}`} />

                {/* Mobile Menu */}
                {mobileOpen && (
                    <div className="md:hidden glass animate-slide-down border-t border-[#1a1a1a]">
                        <div className="px-6 py-6 space-y-4">
                            {navLinks.map((link) => (
                                <a
                                    key={link.label}
                                    href={link.href}
                                    className="block text-sm font-medium text-gray-400 hover:text-white transition-colors"
                                >
                                    {link.label}
                                </a>
                            ))}
                            <div className="pt-4 border-t border-[#1a1a1a]">
                                <a
                                    href="/nget/profile"
                                    className="block text-center text-sm font-semibold px-5 py-2.5 rounded-full border border-gray-700 text-gray-300 hover:border-white hover:text-white transition-all"
                                >
                                    Login
                                </a>
                            </div>
                        </div>
                    </div>
                )}
            </nav>

            {/* Spacer for fixed navbar */}
            <div className="h-[68px]" />
        </>
    )
}

export default Navbar