const fs = require('fs');
const path = require("path");
const vuePlugin = require("esbuild-plugin-vue3");
const tailwindpostcss = require("@tailwindcss/postcss");
const esbuild = require("esbuild");

const BUILD_OPTIONS = {
    logLevel: 'info',
    entryPoints: ['./public/ts/application.ts'],
    bundle: true,
    minify: true,
    minifyWhitespace: true, //avoid leaking routes
    sourcemap: false,
    define: {
        'process.env.NODE_ENV': '"production"',
        'process.env.VERSION': '"1.0.0"'
    },
    // sourcemap: 'inline',
    // sourceRoot: '/js',
    outfile: 'public/js/app.js',
    // target: "browser",
    banner: {
        js: '// javascript',
        css: '/* css */',
    },
    external: ["/img/*", "/fonts/*"],

    conditions: ['style'],
    
    plugins: [
        // Use Vue with Tailwind
        vuePlugin({
            postcss: { plugins: [tailwindpostcss] }
        }),
    ],
    metafile: true,
}

async function build() {
    try {
        let result = await esbuild.build(BUILD_OPTIONS)
        fs.writeFileSync('meta.json', JSON.stringify(result.metafile))
    } catch (e) {
        console.info("Error building, try again...", e)
        // process.exit(1)
    }
}

module.exports = { BUILD_OPTIONS };
watchv2();

async function watchv2() {
    if (process.argv.includes('--watch')) {

        BUILD_OPTIONS.minify = false;
        BUILD_OPTIONS.minifyWhitespace = false;
        BUILD_OPTIONS.sourcemap = true;
        BUILD_OPTIONS.define['process.env.NODE_ENV'] = '"development"';

        await build();

        const watchDirs = [
            path.join(process.cwd(), "pkg/infrastructure/web/templates"),
            path.join(process.cwd(), "public/ts"),
            path.join(process.cwd(), "public/pages"),
        ];

        watchDirs.forEach(dir => {
            fs.watch(dir, { recursive: true }, (eventType, filename) => {
                console.log(`File changed: ${filename} in ${dir}. Rebuilding...`);
                build();
            });
        });
    } else {
        // for some reason this minify is affcting the TailwindCSS v4...
        BUILD_OPTIONS.minify = false;
        // BUILD_OPTIONS.minifyWhitespace = false;

        // build once if not in watch mode
        await build();
        process.exit(0);
    }
}

