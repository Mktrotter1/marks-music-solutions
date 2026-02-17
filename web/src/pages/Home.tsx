import { useRecentAlbums, useRandomAlbums, useStats } from '../api/hooks';
import { AlbumCard } from '../components/AlbumCard';
import { Disc3, Users, Music } from 'lucide-react';

export function HomePage() {
  const { data: recent } = useRecentAlbums(12);
  const { data: random } = useRandomAlbums(12);
  const { data: stats } = useStats();

  return (
    <div className="space-y-8">
      {/* Stats banner */}
      {stats && (stats.artists > 0 || stats.albums > 0 || stats.tracks > 0) && (
        <div className="flex gap-6 text-sm text-mms-text-secondary">
          <div className="flex items-center gap-2">
            <Users size={16} className="text-mms-accent" />
            <span>{stats.artists.toLocaleString()} artists</span>
          </div>
          <div className="flex items-center gap-2">
            <Disc3 size={16} className="text-mms-accent" />
            <span>{stats.albums.toLocaleString()} albums</span>
          </div>
          <div className="flex items-center gap-2">
            <Music size={16} className="text-mms-accent" />
            <span>{stats.tracks.toLocaleString()} tracks</span>
          </div>
        </div>
      )}

      {/* Recently Added */}
      {recent && recent.items.length > 0 && (
        <section>
          <h2 className="text-xl font-bold mb-4">Recently Added</h2>
          <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4">
            {recent.items.map((album) => (
              <AlbumCard key={album.id} album={album} />
            ))}
          </div>
        </section>
      )}

      {/* Discover */}
      {random && random.items.length > 0 && (
        <section>
          <h2 className="text-xl font-bold mb-4">Discover</h2>
          <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4">
            {random.items.map((album) => (
              <AlbumCard key={album.id} album={album} />
            ))}
          </div>
        </section>
      )}

      {/* Empty state */}
      {(!recent || recent.items.length === 0) && (!random || random.items.length === 0) && (
        <div className="text-center py-20">
          <Disc3 size={48} className="mx-auto text-mms-text-tertiary mb-4" />
          <h2 className="text-xl font-semibold mb-2">Your library is empty</h2>
          <p className="text-mms-text-secondary mb-4">
            Configure your music directories in config.yaml, then scan your library.
          </p>
          <p className="text-sm text-mms-text-tertiary">
            Use the "Scan Library" button in the sidebar to get started.
          </p>
        </div>
      )}
    </div>
  );
}
