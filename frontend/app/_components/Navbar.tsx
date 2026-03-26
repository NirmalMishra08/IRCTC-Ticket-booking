import React from 'react'

const Navbar = () => {
    return (
        <div>
            <nav className="bg-white px-6 py-4 flex justify-between items-center border-b">
                <div className="flex items-center gap-2">
                    <div className="w-10 h-10 bg-blue-800 rounded flex items-center justify-center text-white font-bold text-xl">IR</div>
                    <span className="font-extrabold text-2xl text-blue-900 tracking-tighter">IRCTC</span>
                </div>
                <button className="bg-blue-50 text-blue-700 px-6 py-2 rounded-full font-bold text-sm hover:bg-blue-100 transition-all">
                    LOGIN
                </button>
            </nav>
        </div>
    )
}

export default Navbar