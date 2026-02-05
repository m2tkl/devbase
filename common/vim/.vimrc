""
"" Leader Setting
""

let mapleader = "\<Space>"


""
"" Encoding
""

set encoding=utf-8                              " 内部エンコーディング
set fileencodings=utf-8,iso-2022-jp,euc-jp,sjis " 文字コード自動判別(優先順)
set fileformats=unix,dos,mac                    " 改行コード自動判別(優先順)


""
"" Screen
""

set t_Co=256                    " ターミナルで256色表示を使用
set number                      " 行番号の表示
" set ruler                       " 右下に行・列番号を表示
" set list                        " 不可視文字を表示
" 不可視文字の置き換え設定
set listchars=tab:»-,trail:-,extends:»,precedes:«,nbsp:%,eol:↲
set noshowmode                  " 最下部のmode表示をoff (pluginで表示するため)
set showcmd                     " 入力中のコマンドを画面の最下部に表示

syntax on                       " シンタックスハイライト
set cursorline                  " カーソル行をハイライト
" set cursorcolumn                " カーソル位置のカラムをハイライト
set colorcolumn=80              " 80行目に色をつける

set nowrap                      " 折り返さない
set textwidth=0                 " 自動的に改行が入るのを無効化

set background=dark


""
"" System
""

set nocompatible                " Enable arrow key

set nobackup                    " backupファイルを作らない
set nowritebackup               "
set noswapfile                  " swapファイルを作らない

set hidden                      " バッファを閉じる代わりに隠す（Undo履歴を残すため）
set switchbuf=useopen           " 新しく開く代わりにすでに開いてあるバッファを開く
set autoread                    " 外部でファイルが変更された場合は再度読む

set noerrorbells                " エラー音抑制
set novisualbell                " ビジュアルベル抑制
set visualbell t_vb=            " ビープ音抑制

set history=100                 " コマンド、検索履歴数

set scrolloff=5                 " 画面に余裕をもたせてスクロール

set infercase                   " 補完時に大文字小文字を区別しない
set wildmenu                    " ファイル名補完有効
set wildmode=list:longest,full  " 補完モード
" 補完モード description
" - ""                        最初のマッチのみを補完する
" - "full"                    次のマッチを完全に補完する。最後のマッチの次には元の文字列が使われ、その次は再び最初のマッチが補完される
" - "longest"                 共通する最長の文字列までが補完される。それ以上長い文字列を補完できないときは、次の候補に移る
" - "longest:full"            "longest"と似ているが、'wildmenu'が有効ならばそれを開始する。
" - "list"                    複数のマッチがあるときは、すべてのマッチを羅列する。
" - "list:full"               複数のマッチがあるときは、すべてのマッチを羅列し、最初のマッチを補完する。
" - "list:longest"            複数のマッチがあるときは、すべてのマッチを羅列し、共通する最長の文字列までが補完される。

filetype plugin on
syntax on

""
"" Spell check option
""
" Default is off. Space + S to toggle.
""
set spelllang+=cjk              " Japanese is not subject to check spell


""
"" Space + ? to toggle some setting
""
" NOTE: Combine with shift key because toggling is not used very often.
""
nnoremap <silent> <leader><leader>S :setl spell!<CR>:setl spell?<CR>
nnoremap <silent> <leader><leader>L :setl list!<CR>:setl list?<CR>
nnoremap <silent> <leader><leader>T :setl expandtab!<CR>:setl expandtab?<CR>
nnoremap <silent> <leader><leader>W :setl wrap!<CR>:setl wrap?<CR>
nnoremap <silent> <leader><leader>N :setl number!<CR>:setl number?<CR>


""
"" Cursor setting
""
" Options (\e[?)
" 0 : blinking block.
" 1 : blinking block(default).
" 2 : steady block.
" 3 : blinking underline.
" 4 : steady underline.
" 5 : blinking bar(xterm).
" 6 : steady bar(xterm).
""

let &t_SI="\e[5 q"
let &t_EI="\e[2 q"


""
"" Editor setting
""

set clipboard+=unnamed          " OSクリップボード使用
set backspace=indent,eol,start  " バックスペースでなんでも消せるようにする

set tabstop=4                   " <Tab>幅
set softtabstop=4               " <Tab>幅(編集時)
set shiftwidth=4                " インデント幅

set shiftround                  " インデントを'shiftwidth'の倍数に丸める
set expandtab                   " <Tab>押下時に<Tab>ではなく、ホワイトスペースを挿入する
set autoindent                  " 改行時にインデントを引き継いで改行する
set smartindent                 " 前行の末尾に合わせてインデントを増減する

set matchpairs& matchpairs+=<:> " 対応括弧に'<'と'>'のペアを追加
set showmatch                   " 括弧の対応をハイライトする
set matchtime=1                 " 対応表示時間1秒(~10)


augroup fileTypeIndent
    autocmd!
    " python
    autocmd BufNewFile,BufRead *.py setlocal tabstop=4 softtabstop=4 shiftwidth=4
    " elixir
    autocmd BufNewFile,BufRead *.exs setlocal tabstop=2 softtabstop=2 shiftwidth=2
    autocmd BufNewFile,BufRead *.ex setlocal tabstop=2 softtabstop=2 shiftwidth=2
    " frontend
    autocmd BufNewFile,BufRead *.html setlocal tabstop=2 softtabstop=2 shiftwidth=2
    autocmd BufNewFile,BufRead *.css setlocal tabstop=2 softtabstop=2 shiftwidth=2
    autocmd BufNewFile,BufRead *.js setlocal tabstop=2 softtabstop=2 shiftwidth=2
    autocmd BufNewFile,BufRead *.ts setlocal tabstop=2 softtabstop=2 shiftwidth=2
augroup END

augroup autoTagClose
    autocmd!
    autocmd Filetype html inoremap <buffer> ><CR> ></<C-x><C-o><CR><ESC>F>a<CR><ESC>O
    autocmd Filetype html inoremap <buffer> ><Tab> ></<C-x><C-o><CR><ESC>F>a
augroup END

""
"" Search options
""

set ignorecase                  " 大/小文字無視
set smartcase                   " ただし大文字を含む検索は区別
set incsearch                   " インクリメンタルサーチ
set hlsearch                    " 検索マッチテキストをハイライト
set wrapscan                    " 検索をファイルの先頭へループ


""
"" Key mappings
""

" 'jk' to Esc on insert mode
inoremap <silent> jk <ESC>

" Remove highlight by Esc twice
nnoremap <silent> <Esc><Esc> :nohlsearch<CR>

" Search the word under the cursor by *
vnoremap <silent> * "vy/\V<C-r>=substitute(escape(@v, '\/'), "\n", '\\n', 'g')<CR><CR>

" j,k による移動を折り返されたテキストでも自然に振る舞うように変更
nnoremap <silent> j gj
nnoremap <silent> k gk

" Jump bracket pairs
nnoremap <leader>p %
vnoremap <leader>p %

" Split editor (and focus on new editor)
nnoremap <leader><leader>s :<C-u>sp<CR><C-w>j
nnoremap <leader><leader>v :<C-u>vs<CR><C-w>l

" Move to neighboring window
nnoremap <leader><leader>h <C-w>h
nnoremap <leader><leader>j <C-w>j
nnoremap <leader><leader>k <C-w>k
nnoremap <leader><leader>l <C-w>l

" Resize window
nnoremap <silent><silent> <Left>  :vert resize -3<CR>
nnoremap <silent><silent> <Right> :vert resize +3<CR>
nnoremap <silent><silent> <Down>  :resize -3<CR>
nnoremap <silent><silent> <Up>    :resize +3<CR>


" Tab move
nnoremap <Tab>     :tabnext<CR>
nnoremap <S-Tab>   :tabprev<CR>

" <Space>+q closes file.
nnoremap <leader>q :q<CR>

" <Space>+w saves file.
nnoremap <leader>w :w<CR>

" <Space>+1 close file without saving
nnoremap <leader>1 :q!<CR>

" Moves to 行頭、末尾
nnoremap <leader>h ^
vnoremap <leader>h ^
nnoremap <leader>l $
vnoremap <leader>l $

""
"" Terminal Setting
""

" Ctrl+q closes terminal
tnoremap <C-q> <C-w>:q!<CR>


""
"" >>> plugins
""

" Install vim-plug automatically
let data_dir = has('nvim') ? stdpath('data') . '/site' : '~/.vim'
if empty(glob(data_dir . '/autoload/plug.vim'))
      silent execute '!curl -fLo '.data_dir.'/autoload/plug.vim --create-dirs  https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim'
        autocmd VimEnter * PlugInstall --sync | source $MYVIMRC
    endif

" Install plugins automatically
autocmd VimEnter * if len(filter(values(g:plugs), '!isdirectory(v:val.dir)'))
  \| PlugInstall --sync | source $MYVIMRC
  \| endif

" Set plugins
call plug#begin()

" NERDTree
Plug 'scrooloose/nerdtree', { 'on': 'NERDTreeToggle' }

" easymotion
Plug 'easymotion/vim-easymotion'

" lightline
Plug 'itchyny/lightline.vim'

" vim-lsp (https://github.com/mattn/vim-lsp-settings)
Plug 'prabirshrestha/vim-lsp'
Plug 'mattn/vim-lsp-settings'
Plug 'prabirshrestha/asyncomplete.vim'
Plug 'prabirshrestha/asyncomplete-lsp.vim'

" easy-align (https://github.com/junegunn/vim-easy-align)
Plug 'junegunn/vim-easy-align'

" vim-commentary
Plug 'tpope/vim-commentary'

" filer (fern)
Plug 'lambdalisue/fern.vim', { 'on': 'Fern' }

" Emmet-vim (https://github.com/mattn/emmet-vim)
Plug 'mattn/emmet-vim'

" Git
Plug 'airblade/vim-gitgutter'

" ctrlp
Plug 'ctrlpvim/ctrlp.vim'

call plug#end()

""
"" <<< plugins
""

""
"" NERDTree
""

" Show hidden files.
let NERDTreeShowHidden=1

" Set files to be hidden.
let NERDTreeIgnore=['\.DS_Store$','\.localized', 'working.savedSearch']

" Toggle 
map <silent> <leader>b :NERDTreeToggle<CR>


""
"" lightline
""

" Display status bar always.
set laststatus=2


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

" Replace some of the vim seaches with easymotion
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
highlight GitGutterAdd ctermfg=green
highlight GitGutterChange ctermfg=blue
highlight GitGutterDelete ctermfg=red

"" 反映時間を短くする(デフォルトは4000ms)
set updatetime=250
