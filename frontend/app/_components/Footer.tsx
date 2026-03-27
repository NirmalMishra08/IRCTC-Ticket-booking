import React from 'react'

const Footer = () => {
    return (
        <>
            {/* FOOTER */}
            <footer className="mt-auto bg-white py-10 m-2 border-t border-gray-100">
                <div className="max-w-6xl mx-auto px-6 flex flex-col md:flex-row justify-between items-center gap-6">
                    <div className="text-gray-400 text-sm">
                        © 2026 RailGo. An IRCTC Authorized Partner.
                    </div>
                    <div className="flex gap-6 text-sm font-medium text-gray-500">
                        <a href="#" className="hover:text-slate-900">Privacy Policy</a>
                        <a href="#" className="hover:text-slate-900">Terms of Service</a>
                        <a href="#" className="hover:text-slate-900">Contact Us</a>
                    </div>
                </div>
            </footer>
        </>


    )
}

export default Footer