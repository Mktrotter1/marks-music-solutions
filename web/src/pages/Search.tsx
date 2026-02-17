import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { Search as SearchIcon, Disc3, Users, Music } from 'lucide-react';
import { useSearch } from '../api/hooks';

export function SearchPage() {
  const [query, setQuery] = useState('');
  const [debouncedQuery, setDebouncedQuery] = useState('');

  // Debounce search input
  useEffect(() => {
    const timer = setTimeout(() => setDebouncedQuery(query), 300);
    return () => clearTimeout(timer);
  }, [query]);

  const { data: results, isLoading } = useSearch(debouncedQuery);

  return (
    <div>
      {/* Search input */}
      <div className="relative mb-8">
        <SearchIcon
          size={20}
          className="absolute left-4 top-1/2 -translate-y-1/2 text-mms-text-tertiary"
        />
        <input
          data-search-input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Search artists, albums, tracks..."
          className="w-full max-w-[600px] bg-mms-surface border border-mms-border rounded-full pl-12 pr-4 py-3 text-sm text-white placeholder-mms-text-tertiary focus:outline-none focus:border-mms-accent transition-colors"
          autoFocus
        />
      </div>

      {/* Loading */}
      {isLoading && debouncedQuery && (
        <p className="text-mms-text-secondary text-sm">Searching...</p>
      )}

      {/* Results */}
      {results && (
        <div className="space-y-8">
          {/* Artists */}
          {results.artists && results.artists.length > 0 && (
            <section>
              <h3 className="text-lg font-semibold mb-3 flex items-center gap-2">
                <Users size={18} className="text-mms-accent" />
                Artists
              </h3>
              <div className="space-y-1">
                {results.artists.map((r) => (
                  <Link
                    key={r.entity_id}
                    to={`/artist/${r.entity_id}`}
                    className="flex items-center gap-3 px-3 py-2 rounded-md hover:bg-mms-hover transition-colors"
                  >
                    <div className="w-10 h-10 rounded-full bg-mms-surface flex items-center justify-center">
                      <Users size={16} className="text-mms-text-tertiary" />
                    </div>
                    <span className="text-sm">{r.title}</span>
                  </Link>
                ))}
              </div>
            </section>
          )}

          {/* Albums */}
          {results.albums && results.albums.length > 0 && (
            <section>
              <h3 className="text-lg font-semibold mb-3 flex items-center gap-2">
                <Disc3 size={18} className="text-mms-accent" />
                Albums
              </h3>
              <div className="space-y-1">
                {results.albums.map((r) => (
                  <Link
                    key={r.entity_id}
                    to={`/album/${r.entity_id}`}
                    className="flex items-center gap-3 px-3 py-2 rounded-md hover:bg-mms-hover transition-colors"
                  >
                    <div className="w-10 h-10 rounded bg-mms-surface flex items-center justify-center">
                      <Disc3 size={16} className="text-mms-text-tertiary" />
                    </div>
                    <div>
                      <p className="text-sm">{r.title}</p>
                      <p className="text-xs text-mms-text-secondary">{r.artist}</p>
                    </div>
                  </Link>
                ))}
              </div>
            </section>
          )}

          {/* Tracks */}
          {results.tracks && results.tracks.length > 0 && (
            <section>
              <h3 className="text-lg font-semibold mb-3 flex items-center gap-2">
                <Music size={18} className="text-mms-accent" />
                Tracks
              </h3>
              <div className="space-y-1">
                {results.tracks.map((r) => (
                  <div
                    key={r.entity_id}
                    className="flex items-center gap-3 px-3 py-2 rounded-md hover:bg-mms-hover transition-colors"
                  >
                    <div className="w-10 h-10 rounded bg-mms-surface flex items-center justify-center">
                      <Music size={16} className="text-mms-text-tertiary" />
                    </div>
                    <div>
                      <p className="text-sm">{r.title}</p>
                      <p className="text-xs text-mms-text-secondary">
                        {r.artist} &middot; {r.album}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            </section>
          )}

          {/* No results */}
          {results.total === 0 && debouncedQuery && (
            <p className="text-mms-text-secondary text-sm">
              No results found for "{debouncedQuery}"
            </p>
          )}
        </div>
      )}

      {/* Empty state */}
      {!debouncedQuery && (
        <div className="text-center py-16 text-mms-text-tertiary">
          <SearchIcon size={48} className="mx-auto mb-4" />
          <p className="text-lg">Search your music library</p>
          <p className="text-sm mt-1">
            Type to search across artists, albums, and tracks
          </p>
          <p className="text-xs mt-4">
            Tip: Press <kbd className="px-1.5 py-0.5 bg-mms-surface rounded text-mms-text-secondary">/</kbd> from anywhere to jump to search
          </p>
        </div>
      )}
    </div>
  );
}
