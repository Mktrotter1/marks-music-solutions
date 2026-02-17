import { NavLink } from 'react-router-dom';
import {
  Home,
  Search,
  Library,
  Disc3,
  Users,
  ListMusic,
  Settings,
  FolderSync,
} from 'lucide-react';
import { useScanLibrary } from '../api/hooks';

const navItems = [
  { to: '/', icon: Home, label: 'Home' },
  { to: '/search', icon: Search, label: 'Search' },
  { to: '/library', icon: Library, label: 'Library' },
];

const libraryItems = [
  { to: '/library?tab=albums', icon: Disc3, label: 'Albums' },
  { to: '/library?tab=artists', icon: Users, label: 'Artists' },
  { to: '/library?tab=playlists', icon: ListMusic, label: 'Playlists' },
];

export function Sidebar() {
  const scanMutation = useScanLibrary();

  return (
    <div className="flex flex-col h-full py-4">
      {/* Logo */}
      <div className="px-5 mb-6">
        <h1 className="text-lg font-bold tracking-tight">
          <span className="text-mms-accent">MMS</span>
        </h1>
        <p className="text-[10px] text-mms-text-tertiary uppercase tracking-widest">
          Marks Music Server
        </p>
      </div>

      {/* Main navigation */}
      <nav className="px-3 space-y-1">
        {navItems.map(({ to, icon: Icon, label }) => (
          <NavLink
            key={to}
            to={to}
            end={to === '/'}
            className={({ isActive }) =>
              `flex items-center gap-3 px-3 py-2 rounded-md text-sm transition-colors ${
                isActive
                  ? 'bg-mms-surface text-white'
                  : 'text-mms-text-secondary hover:text-white hover:bg-mms-hover'
              }`
            }
          >
            <Icon size={18} />
            {label}
          </NavLink>
        ))}
      </nav>

      {/* Library section */}
      <div className="mt-8 px-3">
        <h2 className="px-3 mb-2 text-[11px] font-semibold uppercase tracking-wider text-mms-text-tertiary">
          Your Library
        </h2>
        <div className="space-y-1">
          {libraryItems.map(({ to, icon: Icon, label }) => (
            <NavLink
              key={to}
              to={to}
              className="flex items-center gap-3 px-3 py-2 rounded-md text-sm text-mms-text-secondary hover:text-white hover:bg-mms-hover transition-colors"
            >
              <Icon size={16} />
              {label}
            </NavLink>
          ))}
        </div>
      </div>

      {/* Spacer */}
      <div className="flex-1" />

      {/* Bottom actions */}
      <div className="px-3 space-y-1">
        <button
          onClick={() => scanMutation.mutate()}
          disabled={scanMutation.isPending}
          className="flex items-center gap-3 px-3 py-2 rounded-md text-sm text-mms-text-secondary hover:text-white hover:bg-mms-hover transition-colors w-full"
        >
          <FolderSync size={16} className={scanMutation.isPending ? 'animate-spin' : ''} />
          {scanMutation.isPending ? 'Scanning...' : 'Scan Library'}
        </button>
        <NavLink
          to="/settings"
          className="flex items-center gap-3 px-3 py-2 rounded-md text-sm text-mms-text-secondary hover:text-white hover:bg-mms-hover transition-colors"
        >
          <Settings size={16} />
          Settings
        </NavLink>
      </div>
    </div>
  );
}
