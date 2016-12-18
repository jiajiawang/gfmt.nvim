if exists('g:loaded_gfmt')
  finish
endif
let g:loaded_gfmt = 1

let g:sort_patterns = {
  \ 'scss': '[a-zA-Z-]+'
\ }

function! s:RequireGfmt(host) abort
  return jobstart(['gfmt'], { 'rpc': v:true })
endfunction

call remote#host#Register('gfmt', 'x', function('s:RequireGfmt'))
call remote#host#RegisterPlugin('gfmt', '0', [
  \ {'type': 'command', 'name': 'Sort', 'sync': 1, 'opts': {'eval': '{''Filetype'': &filetype}', 'nargs': '?', 'range': ''}},
  \ ])

nmap <leader>si vii:Sort<cr>
nmap <leader>sp vip:Sort<cr>
