import { useEffect } from 'react';
import { usePlayerStore } from '../store/player';
import { useNavigate } from 'react-router-dom';

export function useKeyboardShortcuts() {
  const navigate = useNavigate();

  useEffect(() => {
    function handleKeyDown(e: KeyboardEvent) {
      // Don't trigger shortcuts when typing in inputs
      const target = e.target as HTMLElement;
      if (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA') return;

      const player = usePlayerStore.getState();

      switch (e.key) {
        case ' ':
          e.preventDefault();
          player.togglePlay();
          break;
        case 'ArrowRight':
          if (e.shiftKey) {
            player.next();
          } else {
            player.seek(Math.min(player.currentTime + 10, player.duration));
          }
          break;
        case 'ArrowLeft':
          if (e.shiftKey) {
            player.previous();
          } else {
            player.seek(Math.max(player.currentTime - 10, 0));
          }
          break;
        case 'ArrowUp':
          e.preventDefault();
          player.setVolume(Math.min(player.volume + 0.05, 1));
          break;
        case 'ArrowDown':
          e.preventDefault();
          player.setVolume(Math.max(player.volume - 0.05, 0));
          break;
        case 'm':
          player.toggleMute();
          break;
        case 'n':
          player.next();
          break;
        case 'p':
          player.previous();
          break;
        case '/':
          e.preventDefault();
          navigate('/search');
          // Focus search input after navigation
          setTimeout(() => {
            const input = document.querySelector<HTMLInputElement>('[data-search-input]');
            input?.focus();
          }, 100);
          break;
      }
    }

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [navigate]);
}
