""
"" lightline
""

" Display status bar always.
set laststatus=2


""
"" ayu
""

if exists('+termguicolors')
  set termguicolors
endif

let g:ayucolor = 'dark'
colorscheme ayu

function! s:apply_ayu(mode) abort
    let g:ayucolor = a:mode
    execute 'set background=' . a:mode
    colorscheme ayu
endfunction

function! s:toggle_ayu() abort
    if get(g:, 'ayucolor', 'dark') ==# 'light'
        call s:apply_ayu('dark')
    else
        call s:apply_ayu('light')
    endif
endfunction

command! AyuDark call s:apply_ayu('dark')
command! AyuLight call s:apply_ayu('light')
command! AyuToggle call s:toggle_ayu()

nnoremap <silent> <leader><leader>y :AyuToggle<CR>


""
"" easy-align
""

xmap ga <Plug>(EasyAlign)
nmap ga <Plug>(EasyAlign)


""
"" vim-commentary
""
nmap <leader>c gcc
vmap <leader>c gc


""
"" fern:conf
""

" Show hidden files
let g:fern#default_hidden=1

" Show file tree
nnoremap <leader>b :Fern . -drawer -toggle -reveal=%<CR>
nnoremap <leader><leader>e :Fern . -reveal=% <CR>

function! s:init_fern() abort
    " Use 'select' instead of 'edit' for default 'open' action
    nmap <buffer> <Plug>(fern-action-open) <Plug>(fern-action-open:select)
endfunction

augroup fern-custom
    autocmd! *
    autocmd FileType fern call s:init_fern()
augroup END


""
"" easymotion
""
" https://github.com/easymotion/vim-easymotion

let g:EasyMotion_smartcase = 1

" Search in file with 2 chars
nmap <leader>f <Plug>(easymotion-s2)
xmap <leader>f <Plug>(easymotion-s2)

" Replace some of the vim searches with easymotion
map f <Plug>(easymotion-fl)
map F <Plug>(easymotion-Fl)
map t <Plug>(easymotion-tl)
map T <Plug>(easymotion-Tl)


""
"" asyncomplete
""

" Disable completeopt of asyncomplete (Use default vim completeopt)
let g:asyncomplete_auto_completeopt = 0

" No line break with Enter when showing completion candidates
inoremap <expr><CR>  pumvisible() ? "<C-y>" : "<CR>"

" Make the top selected when showing completion.
set completeopt=menuone,noinsert

" C-n, C-p not inserted
inoremap <expr><C-n> pumvisible() ? "<Down>" : "<C-n>"
inoremap <expr><C-p> pumvisible() ? "<Up>" : "<C-p>"


"""
""" Git
"""

" g[で前の変更箇所へ移動する
nnoremap g[ :GitGutterPrevHunk<CR>
" g]で次の変更箇所へ移動する
nnoremap g] :GitGutterNextHunk<CR>
" ghでdiffをハイライトする
nnoremap gh :GitGutterLineHighlightsToggle<CR>
" gpでカーソル行のdiffを表示する
nnoremap gp :GitGutterPreviewHunk<CR>
" 記号の色を変更する
highlight GitGutterAdd ctermfg=green guifg=#aad94c
highlight GitGutterChange ctermfg=blue guifg=#59c2ff
highlight GitGutterDelete ctermfg=red guifg=#f07178

"" 反映時間を短くする(デフォルトは4000ms)
set updatetime=250
