/**
 * printReceipt.js
 * Dedicated thermal receipt printing utility using an isolated iframe.
 * Prevents the main application window and modals from leaking into the browser print dialog.
 */

export function printThermalReceipt(elementOrId, options = {}) {
  const {
    paperWidth = '58mm',
    title = 'Struk Nota'
  } = options

  const targetEl = typeof elementOrId === 'string'
    ? document.getElementById(elementOrId)
    : elementOrId

  if (!targetEl) {
    console.warn(`[printReceipt] Element not found:`, elementOrId)
    return false
  }

  // Remove existing print iframes if any
  const oldIframe = document.getElementById('thermal-print-iframe')
  if (oldIframe) {
    oldIframe.parentNode.removeChild(oldIframe)
  }

  // Create isolated off-screen iframe (must have dimensions > 0 so Chrome won't block it)
  const iframe = document.createElement('iframe')
  iframe.id = 'thermal-print-iframe'
  iframe.style.position = 'fixed'
  iframe.style.top = '-9999px'
  iframe.style.left = '-9999px'
  iframe.style.width = '350px'
  iframe.style.height = '600px'
  iframe.style.border = '0'
  iframe.style.zIndex = '-9999'
  document.body.appendChild(iframe)

  const iframeDoc = iframe.contentWindow.document
  const receiptHtml = targetEl.innerHTML

  const cssWidth = paperWidth === '80mm' ? '80mm' : '58mm'
  const maxWidthPx = paperWidth === '80mm' ? '300px' : '220px'

  iframeDoc.open()
  iframeDoc.write(`
    <!DOCTYPE html>
    <html>
      <head>
        <meta charset="utf-8">
        <title>${title}</title>
        <style>
          @page {
            size: ${cssWidth} auto;
            margin: 0mm;
          }
          *, *::before, *::after {
            box-sizing: border-box;
            margin: 0;
            padding: 0;
          }
          html, body {
            width: ${cssWidth};
            max-width: ${cssWidth};
            margin: 0 auto;
            padding: 6px 4px;
            background: #ffffff !important;
            color: #000000 !important;
            font-family: 'Courier New', Courier, 'Lucida Console', Monaco, monospace;
            font-size: 11px;
            line-height: 1.3;
            -webkit-print-color-adjust: exact;
            print-color-adjust: exact;
          }
          .text-center { text-align: center; }
          .text-left { text-align: left; }
          .text-right { text-align: right; }
          .font-bold { font-weight: bold; }
          .font-extrabold { font-weight: 800; }
          .font-mono { font-family: 'Courier New', Courier, monospace; }
          .uppercase { text-transform: uppercase; }
          .tracking-wider { letter-spacing: 0.05em; }
          .leading-tight { line-height: 1.25; }
          .flex { display: flex; }
          .justify-between { justify-content: space-between; }
          .items-center { align-items: center; }
          .border-b { border-bottom: 1px dashed #000; }
          .border-t { border-top: 1px dashed #000; }
          .border-dashed { border-style: dashed; }
          .border-slate-400, .border-slate-300, .border-slate-200, .border-slate-100 {
            border-color: #000 !important;
          }
          .py-1 { padding-top: 4px; padding-bottom: 4px; }
          .py-1\\.5 { padding-top: 6px; padding-bottom: 6px; }
          .py-2 { padding-top: 8px; padding-bottom: 8px; }
          .py-2\\.5 { padding-top: 10px; padding-bottom: 10px; }
          .py-3 { padding-top: 12px; padding-bottom: 12px; }
          .pb-2 { padding-bottom: 8px; }
          .pb-3 { padding-bottom: 12px; }
          .pt-2 { padding-top: 8px; }
          .pt-3 { padding-top: 12px; }
          .space-y-1 > * + * { margin-top: 4px; }
          .space-y-1\\.5 > * + * { margin-top: 6px; }
          .space-y-2 > * + * { margin-top: 8px; }
          .text-\\[11px\\] { font-size: 11px; }
          .text-\\[10px\\] { font-size: 10px; }
          .text-xs { font-size: 11px; }
          .text-sm { font-size: 12px; }
          .text-base { font-size: 14px; }
          .text-slate-900, .text-slate-800, .text-slate-700, .text-slate-600, .text-slate-500 {
            color: #000000 !important;
          }
          .bg-white, .bg-slate-50 {
            background-color: transparent !important;
          }
          .print\\:hidden, button, svg {
            display: none !important;
          }
        </style>
      </head>
      <body>
        ${receiptHtml}
      </body>
    </html>
  `)
  iframeDoc.close()

  iframe.contentWindow.focus()
  setTimeout(() => {
    try {
      iframe.contentWindow.print()
    } catch (e) {
      console.error('[printReceipt] Error triggering print:', e)
    } finally {
      // Clean up after print dialog finishes
      setTimeout(() => {
        if (iframe && iframe.parentNode) {
          iframe.parentNode.removeChild(iframe)
        }
      }, 1000)
    }
  }, 250)

  return true
}
