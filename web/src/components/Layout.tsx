import { Outlet } from 'react-router-dom';
import { Sidebar } from './Sidebar';
import { NowPlayingBar } from './NowPlayingBar';
import { useKeyboardShortcuts } from '../hooks/useKeyboardShortcuts';
import { usePlayerStore } from '../store/player';

export function Layout() {
  useKeyboardShortcuts();
  const currentTrack = usePlayerStore((s) => s.currentTrack);

  return (
    <div className="h-screen flex flex-col bg-mms-bg overflow-hidden">
      <div className="flex flex-1 overflow-hidden">
        {/* Sidebar */}
        <aside className="w-sidebar flex-shrink-0 bg-mms-sidebar border-r border-mms-border overflow-y-auto hidden md:block">
          <Sidebar />
        </aside>

        {/* Main content */}
        <main className="flex-1 overflow-y-auto">
          <div className="p-6 pb-4">
            <Outlet />
          </div>
        </main>
      </div>

      {/* Now Playing Bar - only show when a track is loaded */}
      {currentTrack && <NowPlayingBar />}
    </div>
  );
}
