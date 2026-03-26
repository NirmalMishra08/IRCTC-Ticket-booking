'use client'
import React, { useState } from 'react';

const TRAIN_SVG = (
  <svg viewBox="0 0 900 220" fill="none" xmlns="http://www.w3.org/2000/svg" className="w-full opacity-10 absolute bottom-0 left-0">
    <rect x="10" y="80" width="880" height="110" rx="18" fill="white"/>
    <rect x="10" y="80" width="160" height="110" rx="18" fill="white"/>
    <circle cx="80" cy="195" r="22" fill="#1e3a8a"/>
    <circle cx="80" cy="195" r="12" fill="white"/>
    <circle cx="200" cy="195" r="22" fill="#1e3a8a"/>
    <circle cx="200" cy="195" r="12" fill="white"/>
    <circle cx="700" cy="195" r="22" fill="#1e3a8a"/>
    <circle cx="700" cy="195" r="12" fill="white"/>
    <circle cx="820" cy="195" r="22" fill="#1e3a8a"/>
    <circle cx="820" cy="195" r="12" fill="white"/>
    <rect x="30" y="100" width="100" height="55" rx="8" fill="#1e3a8a" opacity="0.6"/>
    <rect x="170" y="100" width="120" height="55" rx="6" fill="#1e3a8a" opacity="0.3"/>
    <rect x="310" y="100" width="120" height="55" rx="6" fill="#1e3a8a" opacity="0.3"/>
    <rect x="450" y="100" width="120" height="55" rx="6" fill="#1e3a8a" opacity="0.3"/>
    <rect x="590" y="100" width="120" height="55" rx="6" fill="#1e3a8a" opacity="0.3"/>
    <rect x="730" y="100" width="120" height="55" rx="6" fill="#1e3a8a" opacity="0.3"/>
    <path d="M10 170 Q440 155 890 170" stroke="white" strokeWidth="2" strokeDasharray="8 6"/>
  </svg>
);

export default function IRCTCPage() {
  const [from, setFrom] = useState('New Delhi (NDLS)');
  const [to, setTo] = useState('Mumbai Central (BCT)');
  const [date, setDate] = useState('2026-03-26');
  const [cls, setCls] = useState('all');
  const [tatkal, setTatkal] = useState(true);
  const [flexible, setFlexible] = useState(false);
  const [disability, setDisability] = useState(false);

  const swap = () => { setFrom(to); setTo(from); };

  return (
    <div style={{ fontFamily: "'Sora', 'DM Sans', sans-serif", background: '#f0f4ff', minHeight: '100vh' }}>
      <style>{`
        @import url('https://fonts.googleapis.com/css2?family=Sora:wght@300;400;600;700;800&family=DM+Sans:ital,wght@0,400;0,500;0,700;1,400&display=swap');

        * { box-sizing: border-box; margin: 0; padding: 0; }

        .nav-logo-box {
          width: 44px; height: 44px;
          background: linear-gradient(135deg, #1e3a8a 60%, #3b82f6);
          border-radius: 10px;
          display: flex; align-items: center; justify-content: center;
          font-weight: 900; color: white; font-size: 18px; letter-spacing: -1px;
          box-shadow: 0 4px 12px #1e3a8a44;
        }

        .nav-link {
          font-size: 13px; font-weight: 600; color: #475569;
          text-decoration: none; letter-spacing: 0.02em;
          transition: color 0.2s;
        }
        .nav-link:hover { color: #1e3a8a; }

        .btn-login {
          background: linear-gradient(135deg, #1e3a8a, #2563eb);
          color: white; border: none; padding: 10px 26px;
          border-radius: 50px; font-weight: 700; font-size: 13px;
          cursor: pointer; letter-spacing: 0.06em;
          box-shadow: 0 4px 14px #2563eb44;
          transition: transform 0.15s, box-shadow 0.15s;
        }
        .btn-login:hover { transform: translateY(-1px); box-shadow: 0 6px 20px #2563eb55; }

        /* HERO */
        .hero {
          position: relative;
          background: linear-gradient(160deg, #0f172a 0%, #1e3a8a 50%, #1d4ed8 100%);
          padding: 80px 24px 160px;
          overflow: hidden;
        }
        .hero-grid {
          position: absolute; inset: 0;
          background-image:
            linear-gradient(rgba(255,255,255,0.04) 1px, transparent 1px),
            linear-gradient(90deg, rgba(255,255,255,0.04) 1px, transparent 1px);
          background-size: 48px 48px;
        }
        .hero-glow {
          position: absolute; width: 600px; height: 600px;
          border-radius: 50%;
          background: radial-gradient(circle, #3b82f655 0%, transparent 70%);
          top: -200px; right: -100px; pointer-events: none;
        }
        .hero-glow2 {
          position: absolute; width: 400px; height: 400px;
          border-radius: 50%;
          background: radial-gradient(circle, #f97316 0%, transparent 70%);
          opacity: 0.12;
          bottom: -100px; left: 100px; pointer-events: none;
        }

        .hero-label {
          display: inline-block;
          background: rgba(251,191,36,0.18);
          color: #fbbf24;
          font-size: 11px; font-weight: 700; letter-spacing: 0.15em;
          padding: 5px 14px; border-radius: 50px;
          border: 1px solid rgba(251,191,36,0.3);
          margin-bottom: 18px;
        }
        .hero-title {
          font-size: clamp(32px, 5vw, 54px);
          font-weight: 800; color: white;
          line-height: 1.1; letter-spacing: -0.03em;
          margin-bottom: 12px;
        }
        .hero-title span { color: #fbbf24; }
        .hero-sub {
          font-size: 15px; color: #93c5fd;
          font-weight: 400; max-width: 400px; line-height: 1.6;
        }

        /* SEARCH CARD */
        .search-card {
          background: white;
          border-radius: 24px;
          box-shadow: 0 32px 80px rgba(15,23,42,0.22), 0 0 0 1px rgba(255,255,255,0.6);
          padding: 36px 40px 32px;
          max-width: 860px;
          width: 100%;
        }

        .search-card-title {
          font-size: 13px; font-weight: 700; letter-spacing: 0.1em;
          color: #94a3b8; text-transform: uppercase; margin-bottom: 24px;
        }

        .field-group {
          display: grid;
          grid-template-columns: 1fr auto 1fr 1fr 1fr auto;
          gap: 0;
          border: 1.5px solid #e2e8f0;
          border-radius: 16px;
          overflow: visible;
          position: relative;
          align-items: stretch;
        }

        .field {
          padding: 14px 18px;
          position: relative;
          border-right: 1.5px solid #e2e8f0;
          transition: background 0.15s;
          cursor: text;
        }
        .field:last-child { border-right: none; }
        .field:hover { background: #f8faff; }
        .field:focus-within { background: #eff6ff; }
        .field:focus-within .field-label { color: #2563eb; }

        .field-label {
          font-size: 10px; font-weight: 700; letter-spacing: 0.12em;
          color: #94a3b8; text-transform: uppercase; display: block; margin-bottom: 6px;
          transition: color 0.2s;
        }

        .field-row {
          display: flex; align-items: center; gap: 8px;
        }

        .field-icon { color: #94a3b8; flex-shrink: 0; }

        .field-input {
          border: none; outline: none; background: transparent;
          font-size: 15px; font-weight: 700; color: #1e293b;
          width: 100%; font-family: inherit;
        }
        .field-input::placeholder { color: #cbd5e1; font-weight: 500; }

        .field-select {
          border: none; outline: none; background: transparent;
          font-size: 15px; font-weight: 700; color: #1e293b;
          width: 100%; font-family: inherit; cursor: pointer;
          appearance: none;
        }

        .swap-btn {
          width: 40px; display: flex; align-items: center; justify-content: center;
          border-right: 1.5px solid #e2e8f0; background: transparent; border-left: none;
          cursor: pointer; border-top: none; border-bottom: none;
          flex-shrink: 0;
        }
        .swap-inner {
          width: 32px; height: 32px; background: #eff6ff; border-radius: 50%;
          display: flex; align-items: center; justify-content: center;
          transition: background 0.2s, transform 0.3s;
          border: 1.5px solid #bfdbfe;
        }
        .swap-btn:hover .swap-inner { background: #dbeafe; transform: rotate(180deg); }

        .search-btn {
          background: linear-gradient(135deg, #f97316, #ef4444);
          color: white; border: none; cursor: pointer;
          font-weight: 800; font-size: 14px; letter-spacing: 0.08em;
          border-radius: 0 14px 14px 0;
          padding: 0 28px;
          display: flex; align-items: center; gap: 8px;
          box-shadow: inset 0 1px 0 rgba(255,255,255,0.2);
          transition: opacity 0.2s, transform 0.15s;
          flex-shrink: 0;
        }
        .search-btn:hover { opacity: 0.92; transform: scale(1.01); }

        .options-row {
          margin-top: 20px;
          display: flex; gap: 24px; flex-wrap: wrap; align-items: center;
        }

        .toggle-pill {
          display: flex; align-items: center; gap: 8px;
          font-size: 12px; font-weight: 600; color: #64748b;
          cursor: pointer; user-select: none;
          padding: 6px 14px 6px 8px;
          border-radius: 50px;
          border: 1.5px solid #e2e8f0;
          transition: all 0.2s;
        }
        .toggle-pill.active {
          background: #eff6ff; border-color: #bfdbfe; color: #1d4ed8;
        }
        .toggle-track {
          width: 30px; height: 16px; border-radius: 99px;
          background: #e2e8f0; position: relative;
          transition: background 0.2s; flex-shrink: 0;
        }
        .toggle-track.on { background: #2563eb; }
        .toggle-thumb {
          position: absolute; top: 2px; left: 2px;
          width: 12px; height: 12px; border-radius: 50%;
          background: white; transition: left 0.2s;
          box-shadow: 0 1px 3px rgba(0,0,0,0.2);
        }
        .toggle-track.on .toggle-thumb { left: 16px; }

        /* CARDS SECTION */
        .cards-wrap {
          max-width: 900px; margin: 0 auto;
          padding: 0 24px;
          margin-top: -60px;
          position: relative; z-index: 10;
          display: grid; grid-template-columns: repeat(auto-fit, minmax(220px, 1fr)); gap: 20px;
        }

        .info-card {
          background: white;
          border-radius: 20px;
          padding: 28px 24px;
          box-shadow: 0 8px 32px rgba(15,23,42,0.08);
          border: 1px solid rgba(226,232,240,0.8);
          display: flex; align-items: flex-start; gap: 16px;
          transition: transform 0.2s, box-shadow 0.2s;
        }
        .info-card:hover { transform: translateY(-4px); box-shadow: 0 16px 48px rgba(15,23,42,0.12); }

        .info-icon {
          width: 48px; height: 48px; border-radius: 14px;
          display: flex; align-items: center; justify-content: center;
          font-size: 22px; flex-shrink: 0;
        }

        .info-card-title {
          font-size: 14px; font-weight: 700; color: #1e293b; margin-bottom: 4px;
        }
        .info-card-desc {
          font-size: 12px; color: #94a3b8; line-height: 1.5; font-weight: 400;
        }

        /* QUICK LINKS */
        .quick-section {
          max-width: 900px; margin: 60px auto 0;
          padding: 0 24px 80px;
        }
        .section-label {
          font-size: 11px; font-weight: 700; letter-spacing: 0.14em;
          color: #94a3b8; text-transform: uppercase; margin-bottom: 20px;
        }
        .quick-grid {
          display: grid; grid-template-columns: repeat(auto-fill, minmax(130px, 1fr)); gap: 12px;
        }
        .quick-tile {
          border: 1.5px solid #e2e8f0; border-radius: 16px;
          padding: 18px 14px; text-align: center;
          cursor: pointer; transition: all 0.2s;
          background: white;
        }
        .quick-tile:hover { border-color: #93c5fd; background: #eff6ff; transform: translateY(-2px); }
        .quick-tile-icon { font-size: 26px; margin-bottom: 8px; }
        .quick-tile-label { font-size: 11px; font-weight: 700; color: #475569; letter-spacing: 0.02em; }

        /* FOOTER */
        .footer {
          background: #0f172a; color: #475569;
          text-align: center; padding: 24px;
          font-size: 12px; letter-spacing: 0.04em;
        }
        .footer span { color: #64748b; }

        @media (max-width: 720px) {
          .field-group { grid-template-columns: 1fr; }
          .search-btn { border-radius: 0 0 14px 14px; padding: 14px; justify-content: center; }
          .swap-btn { display: none; }
          .field { border-right: none; border-bottom: 1.5px solid #e2e8f0; }
          .hero { padding: 60px 20px 140px; }
          .search-card { padding: 24px 20px; }
          .hero-title { font-size: 30px; }
        }
      `}</style>

      {/* NAV */}
      <nav style={{
        background: 'rgba(15,23,42,0.97)',
        backdropFilter: 'blur(16px)',
        padding: '0 32px',
        display: 'flex', alignItems: 'center', justifyContent: 'space-between',
        height: '64px',
        position: 'sticky', top: 0, zIndex: 100,
        borderBottom: '1px solid rgba(255,255,255,0.06)'
      }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
          <div className="nav-logo-box">IR</div>
          <div>
            <div style={{ fontWeight: 800, fontSize: '16px', color: 'white', letterSpacing: '-0.04em', lineHeight: 1 }}>IRCTC</div>
            <div style={{ fontSize: '9px', color: '#64748b', letterSpacing: '0.14em', fontWeight: 600 }}>INDIAN RAILWAYS</div>
          </div>
        </div>

        <div style={{ display: 'flex', gap: '28px', alignItems: 'center' }}>
          {['Trains', 'PNR Status', 'Fare Alert', 'My Trips'].map(l => (
            <a key={l} className="nav-link" href="#" style={{ color: '#94a3b8' }}>{l}</a>
          ))}
          <button className="btn-login">SIGN IN</button>
        </div>
      </nav>

      {/* HERO */}
      <section className="hero">
        <div className="hero-grid" />
        <div className="hero-glow" />
        <div className="hero-glow2" />
        {TRAIN_SVG}

        <div style={{ maxWidth: '1100px', margin: '0 auto', position: 'relative', zIndex: 2, display: 'flex', gap: '80px', alignItems: 'flex-start', flexWrap: 'wrap' }}>
          {/* Left copy */}
          <div style={{ paddingTop: '12px', flex: '0 0 300px' }}>
            <div className="hero-label">🚆 BOOKING OPEN</div>
            <h1 className="hero-title">
              Book your<br /><span>journey</span><br />across India.
            </h1>
            <p className="hero-sub">
              Fast, reliable, and secure ticket booking for 13,000+ trains to 7,000+ stations.
            </p>
            <div style={{ marginTop: '28px', display: 'flex', gap: '24px' }}>
              {[['14K+', 'Trains'], ['7K+', 'Stations'], ['1M+', 'Daily Users']].map(([n, l]) => (
                <div key={l}>
                  <div style={{ fontWeight: 800, fontSize: '22px', color: 'white' }}>{n}</div>
                  <div style={{ fontSize: '11px', color: '#64748b', fontWeight: 600, letterSpacing: '0.08em' }}>{l}</div>
                </div>
              ))}
            </div>
          </div>

          {/* Search card */}
          <div className="search-card" style={{ flex: 1, minWidth: '320px' }}>
            <div className="search-card-title">Where are you headed?</div>

            <div className="field-group">
              {/* From */}
              <div className="field" style={{ borderRadius: '14px 0 0 14px' }}>
                <span className="field-label">From</span>
                <div className="field-row">
                  <svg className="field-icon" width="16" height="16" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}><circle cx="12" cy="10" r="3"/><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7z"/></svg>
                  <input className="field-input" value={from} onChange={e => setFrom(e.target.value)} placeholder="Origin station" />
                </div>
              </div>

              {/* Swap */}
              <button className="swap-btn" onClick={swap}>
                <div className="swap-inner">
                  <svg width="14" height="14" fill="none" viewBox="0 0 24 24" stroke="#2563eb" strokeWidth={2.5}><path d="M7 16V4m0 0L3 8m4-4l4 4M17 8v12m0 0l4-4m-4 4l-4-4"/></svg>
                </div>
              </button>

              {/* To */}
              <div className="field">
                <span className="field-label">To</span>
                <div className="field-row">
                  <svg className="field-icon" width="16" height="16" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7z"/><circle cx="12" cy="10" r="3" fill="currentColor" stroke="none"/></svg>
                  <input className="field-input" value={to} onChange={e => setTo(e.target.value)} placeholder="Destination" />
                </div>
              </div>

              {/* Date */}
              <div className="field">
                <span className="field-label">Date</span>
                <div className="field-row">
                  <svg className="field-icon" width="16" height="16" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}><rect x="3" y="4" width="18" height="18" rx="2"/><path d="M16 2v4M8 2v4M3 10h18"/></svg>
                  <input type="date" className="field-input" value={date} onChange={e => setDate(e.target.value)} style={{ fontSize: '13px' }} />
                </div>
              </div>

              {/* Class */}
              <div className="field">
                <span className="field-label">Class</span>
                <div className="field-row">
                  <svg className="field-icon" width="16" height="16" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}><rect x="2" y="7" width="20" height="14" rx="2"/><path d="M16 7V5a2 2 0 00-2-2h-4a2 2 0 00-2 2v2"/></svg>
                  <select className="field-select" value={cls} onChange={e => setCls(e.target.value)}>
                    <option value="all">All Classes</option>
                    <option value="sl">Sleeper (SL)</option>
                    <option value="3a">3 Tier AC (3A)</option>
                    <option value="2a">2 Tier AC (2A)</option>
                    <option value="1a">First AC (1A)</option>
                    <option value="cc">Chair Car (CC)</option>
                  </select>
                </div>
              </div>

              {/* Search Btn */}
              <button className="search-btn">
                <svg width="18" height="18" fill="none" viewBox="0 0 24 24" stroke="white" strokeWidth={2.5}><circle cx="11" cy="11" r="8"/><path d="m21 21-4.35-4.35"/></svg>
                SEARCH
              </button>
            </div>

            {/* Toggles */}
            <div className="options-row">
              {[
                [tatkal, setTatkal, '⚡', 'Tatkal'],
                [flexible, setFlexible, '📅', 'Flexible Date'],
                [disability, setDisability, '♿', 'Divyaang'],
              ].map(([val, set, icon, label]) => (
                <label key={label} className={`toggle-pill ${val ? 'active' : ''}`} onClick={() => set(!val)}>
                  <div className={`toggle-track ${val ? 'on' : ''}`}>
                    <div className="toggle-thumb" />
                  </div>
                  <span>{icon} {label}</span>
                </label>
              ))}
            </div>
          </div>
        </div>
      </section>

      {/* INFO CARDS */}
      <div className="cards-wrap">
        {[
          { icon: '⚡', color: '#fef9c3', label: 'Instant Booking', desc: 'Confirm your seat in under 60 seconds.', accent: '#eab308' },
          { icon: '🔒', color: '#dcfce7', label: 'Secure Payments', desc: '100% safe with bank-grade encryption.', accent: '#16a34a' },
          { icon: '💸', color: '#dbeafe', label: 'Quick Refunds', desc: 'Refunds directly to your source account.', accent: '#2563eb' },
          { icon: '📱', color: '#fce7f3', label: '24/7 Support', desc: 'Help center available round-the-clock.', accent: '#db2777' },
        ].map(c => (
          <div className="info-card" key={c.label}>
            <div className="info-icon" style={{ background: c.color }}>
              {c.icon}
            </div>
            <div>
              <div className="info-card-title">{c.label}</div>
              <div className="info-card-desc">{c.desc}</div>
            </div>
          </div>
        ))}
      </div>

      {/* QUICK LINKS */}
      <div className="quick-section">
        <div className="section-label">Quick Services</div>
        <div className="quick-grid">
          {[
            ['🎫', 'Book Ticket'],
            ['📋', 'PNR Status'],
            ['🔔', 'Fare Alert'],
            ['🗺️', 'Train Route'],
            ['⏱️', 'Live Status'],
            ['♻️', 'Cancel Ticket'],
            ['🏨', 'Rail Hotel'],
            ['🍱', 'Order Food'],
          ].map(([icon, label]) => (
            <div className="quick-tile" key={label}>
              <div className="quick-tile-icon">{icon}</div>
              <div className="quick-tile-label">{label}</div>
            </div>
          ))}
        </div>
      </div>

      {/* FOOTER */}
      <div className="footer">
        <span>© 2026 Indian Railway Catering and Tourism Corporation Ltd.</span>
        &nbsp;·&nbsp; A Government of India Enterprise
      </div>
    </div>
  );
}