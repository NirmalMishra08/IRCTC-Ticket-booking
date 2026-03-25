import React from 'react';
import { Search, ArrowRightLeft, Calendar, MapPin } from 'lucide-react';

export default function Page() {
  return (
    <div className="min-h-screen bg-slate-50">
      {/* --- MINIMALIST NAV --- */}
      <nav className="bg-white px-6 py-4 flex justify-between items-center border-b">
        <div className="flex items-center gap-2">
          <div className="w-10 h-10 bg-blue-800 rounded flex items-center justify-center text-white font-bold text-xl">IR</div>
          <span className="font-extrabold text-2xl text-blue-900 tracking-tighter">IRCTC</span>
        </div>
        <button className="bg-blue-50 text-blue-700 px-6 py-2 rounded-full font-bold text-sm hover:bg-blue-100 transition-all">
          LOGIN
        </button>
      </nav>

      {/* --- HERO SECTION WITH IMAGE & SEARCH --- */}
      <section className="relative h-[500px] flex items-center justify-center px-4">
        {/* Background Image with Overlay */}
        <div 
          className="absolute inset-0 bg-cover bg-center z-0" 
          style={{ backgroundImage: `url('/train.jpg')` }}
        >
          <div className="absolute inset-0 bg-blue-900/40 backdrop-blur-[2px]"></div>
        </div>

        {/* --- COMPACT SEARCH BOX --- */}
        <div className="relative z-10 w-full max-w-5xl bg-white/95 backdrop-blur-md rounded-2xl shadow-2xl p-8">
          <h1 className="text-2xl font-bold text-slate-800 mb-6 flex items-center gap-2">
            BOOK TICKET
          </h1>

          <div className="grid grid-cols-1 md:grid-cols-12 gap-4 items-end">
            {/* From & To with Swap */}
            <div className="md:col-span-5 relative grid grid-cols-2 gap-0 border rounded-xl overflow-hidden shadow-sm">
              <div className="p-3 border-r">
                <label className="text-[10px] font-bold text-blue-600 uppercase block mb-1">From</label>
                <div className="flex items-center gap-2">
                  <MapPin size={16} className="text-slate-400" />
                  <input type="text" placeholder="Origin" className="w-full outline-none font-bold text-slate-700" defaultValue="NDLS" />
                </div>
              </div>
              
              {/* Absolute Swap Button */}
              <button className="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 w-8 h-8 bg-white border shadow-md rounded-full flex items-center justify-center z-10 hover:bg-slate-50 transition-colors">
                <ArrowRightLeft size={14} className="text-blue-600" />
              </button>

              <div className="p-3 pl-6">
                <label className="text-[10px] font-bold text-blue-600 uppercase block mb-1 text-right">To</label>
                <div className="flex items-center gap-2">
                  <input type="text" placeholder="Destination" className="w-full outline-none font-bold text-slate-700 text-right" defaultValue="BCT" />
                  <MapPin size={16} className="text-slate-400" />
                </div>
              </div>
            </div>

            {/* Date Picker */}
            <div className="md:col-span-3 border rounded-xl p-3 shadow-sm">
              <label className="text-[10px] font-bold text-blue-600 uppercase block mb-1">Date of Journey</label>
              <div className="flex items-center gap-2">
                <Calendar size={18} className="text-slate-400" />
                <input type="text" defaultValue="26 Mar 2026" className="w-full outline-none font-bold text-slate-700" />
              </div>
            </div>

            {/* Class Dropdown */}
            <div className="md:col-span-2 border rounded-xl p-3 shadow-sm">
              <label className="text-[10px] font-bold text-blue-600 uppercase block mb-1">Class</label>
              <select className="w-full outline-none font-bold text-slate-700 bg-transparent cursor-pointer">
                <option>All Classes</option>
                <option>Sleeper</option>
                <option>3 Tier AC</option>
              </select>
            </div>

            {/* Search Button */}
            <div className="md:col-span-2">
              <button className="w-full bg-orange-500 hover:bg-orange-600 text-white font-bold h-[58px] rounded-xl shadow-lg shadow-orange-200 transition-all flex items-center justify-center gap-2">
                <Search size={20} /> SEARCH
              </button>
            </div>
          </div>

          {/* Quick Shortcuts */}
          <div className="mt-6 flex flex-wrap gap-4 text-xs font-bold text-slate-500">
            <label className="flex items-center gap-2 cursor-pointer hover:text-blue-600">
              <input type="checkbox" className="accent-blue-600" /> PERSON WITH DISABILITY
            </label>
            <label className="flex items-center gap-2 cursor-pointer hover:text-blue-600">
              <input type="checkbox" className="accent-blue-600" /> FLEXIBLE WITH DATE
            </label>
            <label className="flex items-center gap-2 cursor-pointer hover:text-blue-600">
              <input type="checkbox" className="accent-blue-600 text-white" checked /> TATKAL
            </label>
          </div>
        </div>
      </section>

      {/* --- INFO TILES --- */}
      <div className="max-w-5xl mx-auto -translate-y-10 grid grid-cols-1 md:grid-cols-3 gap-6 px-4">
        <InfoCard title="Easy Booking" desc="Book your tickets in under 2 minutes." icon="⚡" />
        <InfoCard title="Instant Refund" desc="Get refunds directly to your source bank." icon="💰" />
        <InfoCard title="24/7 Support" desc="We are here to help you anytime." icon="📞" />
      </div>
    </div>
  );
}

function InfoCard({ title, desc, icon }: { title: string, desc: string, icon: string }) {
  return (
    <div className="bg-white p-6 rounded-xl shadow-xl border-t-4 border-blue-600 flex items-center gap-4">
      <span className="text-3xl">{icon}</span>
      <div>
        <h4 className="font-bold text-slate-800">{title}</h4>
        <p className="text-xs text-slate-500">{desc}</p>
      </div>
    </div>
  );
}