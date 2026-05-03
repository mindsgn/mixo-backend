package admin

const adminPageTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Radio Admin</title>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            background: #f5f5f5;
            padding: 20px;
        }
        .container { max-width: 1200px; margin: 0 auto; }
        h1 { color: #333; margin-bottom: 30px; }
        .nav { margin-bottom: 30px; }
        .nav a { color: #667eea; text-decoration: none; margin-right: 20px; font-weight: 500; }
        .nav a:hover { text-decoration: underline; }
        .section {
            background: white;
            border-radius: 10px;
            padding: 30px;
            margin-bottom: 30px;
            box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
        }
        h2 { color: #333; margin-bottom: 20px; font-size: 22px; }
        .form-group { margin-bottom: 15px; }
        label { display: block; margin-bottom: 5px; color: #555; font-weight: 500; }
        input, select {
            width: 100%;
            padding: 10px;
            border: 1px solid #ddd;
            border-radius: 5px;
            font-size: 14px;
        }
        input:focus, select:focus { outline: none; border-color: #667eea; }
        button {
            background: #667eea;
            color: white;
            border: none;
            padding: 10px 20px;
            border-radius: 5px;
            cursor: pointer;
            font-size: 14px;
            font-weight: 500;
        }
        button:hover { background: #5568d3; }
        button.delete { background: #dc3545; }
        button.delete:hover { background: #c82333; }
        button.add-queue { background: #28a745; }
        button.add-queue:hover { background: #218838; }
        button.play { background: #17a2b8; }
        button.play:hover { background: #138496; }
        button.stop { background: #ffc107; color: #333; }
        button.stop:hover { background: #e0a800; }
        table { width: 100%%; border-collapse: collapse; margin-top: 20px; }
        th, td { padding: 12px; text-align: left; border-bottom: 1px solid #eee; }
        th { background: #f8f9fa; color: #555; font-weight: 600; }
        tr:hover { background: #f8f9fa; }
        .actions { display: flex; gap: 10px; }
        .empty { color: #999; font-style: italic; padding: 20px; text-align: center; }
        .error { background: #f8d7da; color: #721c24; padding: 10px; border-radius: 5px; margin-bottom: 15px; }
        .success { background: #d4edda; color: #155724; padding: 10px; border-radius: 5px; margin-bottom: 15px; }
        .now-playing {
            background: #e3f2fd;
            border-left: 4px solid #2196f3;
            padding: 15px;
            margin-bottom: 20px;
            border-radius: 5px;
        }
        .now-playing h3 { margin: 0 0 10px 0; color: #1976d2; }
        .now-playing p { margin: 5px 0; color: #555; }
        .status { font-weight: bold; }
        .status.playing { color: #28a745; }
        .status.paused { color: #ffc107; }
        [hx-indicator] { opacity: 0.5; pointer-events: none; }
    </style>
</head>
<body>
    <div class="container">
        <div class="nav">
            <a href="/stream">Listen to Stream</a>
        </div>
        <h1>Radio Admin</h1>
        <div id="message"></div>

        <div class="section">
            <h2>Playback Control</h2>
            <div id="now-playing" hx-get="/admin/now-playing" hx-trigger="load, every 5s">
                %s
            </div>
        </div>

        <div class="section">
            <h2>Upload MP3 File</h2>
            <form hx-post="/admin/upload" hx-encoding="multipart/form-data" hx-target="#message" hx-swap="innerHTML" hx-on::after-request="if(event.detail.successful) { this.reset(); htmx.trigger('#songs-table', 'refresh'); }">
                <div class="form-group">
                    <label for="file">MP3 File</label>
                    <input type="file" id="file" name="file" accept=".mp3,audio/mpeg" required>
                </div>
                <button type="submit">Upload Song</button>
            </form>
        </div>

        <div class="section">
            <h2>Songs Library</h2>
            <div id="songs-table" hx-get="/admin/songs" hx-trigger="load, refresh">
                %s
            </div>
        </div>

        <div class="section">
            <h2>Playback Queue</h2>
            <div id="queue-table" hx-get="/admin/queue" hx-trigger="load, refresh, every 10s">
                %s
            </div>
        </div>
    </div>
</body>
</html>`

const songsTableTemplate = `<table>
    <thead>
        <tr>
            <th>Title</th>
            <th>Artist</th>
            <th>Duration</th>
            <th>Actions</th>
        </tr>
    </thead>
    <tbody>
        %s
    </tbody>
</table>`

const songRowTemplate = `<tr>
    <td>%s</td>
    <td>%s</td>
    <td>%ds</td>
    <td class="actions">
        <button class="add-queue" hx-post="/admin/queue/%d" hx-target="#message" hx-swap="innerHTML" hx-on::after-request="htmx.trigger('#queue-table', 'refresh')">Add to Queue</button>
        <button class="delete" hx-delete="/admin/songs/%d" hx-target="#message" hx-swap="innerHTML" hx-confirm="Delete this song?" hx-on::after-request="htmx.trigger('#songs-table', 'refresh')">Delete</button>
    </td>
</tr>`

const emptySongsTemplate = `<div class="empty">No songs in library</div>`

const queueTableTemplate = `<table>
    <thead>
        <tr>
            <th>Position</th>
            <th>Title</th>
            <th>Artist</th>
            <th>Duration</th>
            <th>Actions</th>
        </tr>
    </thead>
    <tbody>
        %s
    </tbody>
</table>`

const queueRowTemplate = `<tr>
    <td>%d</td>
    <td>%s</td>
    <td>%s</td>
    <td>%ds</td>
    <td class="actions">
        <button class="delete" hx-delete="/admin/queue/%d" hx-target="#message" hx-swap="innerHTML" hx-on::after-request="htmx.trigger('#queue-table', 'refresh')">Remove</button>
    </td>
</tr>`

const emptyQueueTemplate = `<div class="empty">Queue is empty</div>`

const nowPlayingTemplate = `<div class="now-playing">
    <h3>Now Playing</h3>
    <p><strong>%s</strong> by %s</p>
    <p>Duration: %ds</p>
    <p>Status: <span class="status %s">%s</span></p>
    <form hx-post="/admin/play" hx-target="#now-playing" hx-swap="innerHTML">
        <button type="submit" class="%s">%s</button>
    </form>
</div>`

const nowPlayingEmptyTemplate = `<div class="now-playing">
    <h3>Now Playing</h3>
    <p>No song currently playing</p>
    <p>Status: <span class="status %s">%s</span></p>
    <form hx-post="/admin/play" hx-target="#now-playing" hx-swap="innerHTML">
        <button type="submit" class="%s">%s</button>
    </form>
</div>`

const messageSuccessTemplate = `<div class="success">%s</div>`
const messageErrorTemplate = `<div class="error">%s</div>`
