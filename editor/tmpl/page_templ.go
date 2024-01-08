// Code generated by templ - DO NOT EDIT.

// templ: version: 0.2.476
package tmpl

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func page(title string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!doctype html><html lang=\"en\"><head><title>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string = title
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</title><meta charset=\"utf-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"><link rel=\"stylesheet\" href=\"https://unpkg.com/primeflex@3.3.1/primeflex.css\"><link rel=\"stylesheet\" href=\"https://unpkg.com/primeflex@3.3.1/themes/primeone-light.css\"><link rel=\"stylesheet\" href=\"https://unpkg.com/primeicons@4.1.0/primeicons.css\"><script src=\"https://unpkg.com/slim-select@latest/dist/slimselect.min.js\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Var3 := ``
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var3)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</script><link href=\"https://unpkg.com/slim-select@latest/dist/slimselect.css\" rel=\"stylesheet\"><style>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Var4 := `
                .container {
                    max-width: 1200px;
                    margin: 0 auto;
                }
                .lnt-tag {
                    display: inline-block;
                    padding: 0.25rem 0.5rem;
                    border-radius: 0.25rem;
                    margin: 0 0.25rem 0.25rem 0;
                    font-size: 0.95rem;
                    font-weight: 700;
                    font-family: 'Noto Sans JP', sans-serif;
                    background-color: var(--primary-700);
                    white-space: nowrap;
                    text-align: center;
                }
                .lnt-tag-sm {
                    font-size: 0.82rem;
                    padding: 0.15rem 0.3rem;
                }
                .bg-light-novel, .bg-novel {
                    background-color: var(--cyan-600);
                }
                .bg-manga {
                    background-color: var(--purple-600);
                }
                .bg-artbook {
                    background-color: var(--orange-600);
                }
                .series-list .odd {
                    background-color: var(--gray-200);
                }
                .m-h-2rem {
                    min-height: 2rem;
                }
                .field {
                    padding-left: 0.5rem;
                    padding-right: 0.5rem;
                }
                .field label {
                    width: 100%;
                    font-size: 0.9rem;
                }
                .field input {
                    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol";
                    font-size: 1rem;
                    color: #495057;
                    background: #ffffff;
                    padding: 0.5rem 0.5rem;
                    border: 1px solid #ced4da;
                    transition: background-color 0.2s, color 0.2s, border-color 0.2s, box-shadow 0.2s;
                    appearance: none;
                    border-radius: 6px;
                    width: 100%;
                }
                .field select {
                    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol";
                    font-size: 1rem;
                    color: #495057;
                    background: #ffffff;
                    padding: 0.5rem 0.5rem;
                    border: 1px solid #ced4da;
                    transition: background-color 0.2s, color 0.2s, border-color 0.2s, box-shadow 0.2s;
                    border-radius: 6px;
                    width: 100%;
                }
                .field textarea {
                    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol";
                    font-size: 1rem;
                    color: #495057;
                    background: #ffffff;
                    padding: 0.5rem 0.5rem;
                    border: 1px solid #ced4da;
                    transition: background-color 0.2s, color 0.2s, border-color 0.2s, box-shadow 0.2s;
                    appearance: none;
                    border-radius: 6px;
                    width: 100%;
                }
                .btn {
                    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol";
                    font-size: 1rem;
                    color: #ffffff;
                    background: #007bff;
                    padding: 0.3rem 0.6rem;
                    border: 1px solid #007bff;
                    transition: background-color 0.2s, color 0.2s, border-color 0.2s, box-shadow 0.2s;
                    appearance: none;
                    border-radius: 6px;
                    cursor: pointer;
                    text-decoration: none;
                    display: inline-block;
                    text-align: center;
                    vertical-align: middle;
                    user-select: none;
                    white-space: nowrap;
                    line-height: 1.5;
                    font-size: 1rem;
                    border-radius: 0.25rem;
                    font-weight: 700;
                }
                .btn-primary {
                    background-color: var(--primary-500);
                    border-color: var(--primary-500);
                }
                .btn-secondary {
                    background-color: var(--secondary-500);
                    border-color: var(--secondary-500);
                }
                .btn-danger {
                    background-color: var(--red-500);
                    border-color: var(--red-500);
                }
            `
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var4)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</style></head><body><div class=\"container\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ_7745c5c3_Var1.Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
